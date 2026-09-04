package utils

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Manager struct {
	PublicURL string

	cmd    *exec.Cmd
	cancel context.CancelFunc

	mu      sync.RWMutex
	started bool
}

var GlobalTunnel = &Manager{}

var cloudflareURLRegex = regexp.MustCompile(
	`https://[a-zA-Z0-9-]+\.trycloudflare\.com`,
)


const maxStderrLines = 5

func (tm *Manager) StartCloudflareTunnel(ctx context.Context) (string, error) {
	tm.mu.Lock()

	if tm.started && tm.PublicURL != "" {
		url := tm.PublicURL
		tm.mu.Unlock()
		return url, nil
	}

	tm.mu.Unlock()

	cloudflaredPath := os.Getenv("CLOUDFLARED_PATH")

	if cloudflaredPath == "" {
		return "", fmt.Errorf(
			"CLOUDFLARED_PATH environment variable is not set",
		)
	}

	if _, err := os.Stat(cloudflaredPath); err != nil {
		return "", fmt.Errorf(
			"cloudflared executable not found at %s: %w",
			cloudflaredPath,
			err,
		)
	}

	procCtx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(
		procCtx,
		cloudflaredPath,
		"tunnel",
		"--url",
		"http://127.0.0.1:8080",
		"--edge-ip-version",
		"4",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return "", fmt.Errorf(
			"failed to create cloudflared stdout pipe: %w",
			err,
		)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		return "", fmt.Errorf(
			"failed to create cloudflared stderr pipe: %w",
			err,
		)
	}

	if err := cmd.Start(); err != nil {
		cancel()
		return "", fmt.Errorf(
			"failed to start cloudflared: %w",
			err,
		)
	}

	tm.mu.Lock()
	tm.cmd = cmd
	tm.cancel = cancel
	tm.mu.Unlock()

	urlChan := make(chan string, 1)
	errChan := make(chan error, 1)

	var (
		lastLinesMu sync.Mutex
		lastLines   []string
	)

	appendLine := func(line string) {
		lastLinesMu.Lock()
		defer lastLinesMu.Unlock()

		lastLines = append(lastLines, line)
		if len(lastLines) > maxStderrLines {
			lastLines = lastLines[len(lastLines)-maxStderrLines:]
		}
	}

	snapshotLines := func() string {
		lastLinesMu.Lock()
		defer lastLinesMu.Unlock()

		return strings.Join(lastLines, " | ")
	}


	go func() {
		scanner := bufio.NewScanner(stderr)

		for scanner.Scan() {
			line := scanner.Text()

			log.Printf("[Cloudflare] %s", line)
			appendLine(line)

			match := cloudflareURLRegex.FindString(line)

			if match != "" {
				if strings.HasPrefix(match, "https://api.trycloudflare.com") {
					continue
				}

				select {
				case urlChan <- match:
				default:
				}

				continue
			}
		}

		if err := scanner.Err(); err != nil {
			select {
			case errChan <- fmt.Errorf(
				"failed reading cloudflared output: %w (last output: %s)",
				err,
				snapshotLines(),
			):
			default:
			}

			return
		}

		select {
		case errChan <- fmt.Errorf(
			"cloudflared exited before generating a public URL, last output: %s",
			snapshotLines(),
		):
		default:
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stdout)

		for scanner.Scan() {
			log.Printf("[Cloudflare] %s", scanner.Text())
		}
	}()

	// Monitor the process.
	go func() {
		err := cmd.Wait()

		tm.mu.Lock()
		tm.started = false
		tm.PublicURL = ""
		tm.cmd = nil
		tm.cancel = nil
		tm.mu.Unlock()

		if err != nil {
			log.Printf("[Cloudflare] tunnel stopped: %v", err)
		} else {
			log.Println("[Cloudflare] tunnel stopped")
		}
	}()

	select {
	case publicURL := <-urlChan:

		publicURL = strings.TrimRight(
			publicURL,
			"/",
		)

		tm.mu.Lock()
		tm.PublicURL = publicURL
		tm.started = true
		tm.mu.Unlock()

		log.Printf(
			"[Cloudflare] Public URL: %s",
			publicURL,
		)

		return publicURL, nil

	case err := <-errChan:

		cancel()

		return "", err

	case <-ctx.Done():

		cancel()

		return "", fmt.Errorf(
			"timed out waiting for Cloudflare tunnel: %w",
			ctx.Err(),
		)
	}
}


func (tm *Manager) StartCloudflareTunnelWithRetry(
	ctx context.Context,
	maxAttempts int,
	perAttemptTimeout time.Duration,
) (string, error) {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		attemptCtx, cancel := context.WithTimeout(ctx, perAttemptTimeout)

		url, err := tm.StartCloudflareTunnel(attemptCtx)

		cancel()

		if err == nil {
			return url, nil
		}

		lastErr = err

		log.Printf(
			"[Cloudflare] attempt %d/%d failed: %v",
			attempt,
			maxAttempts,
			err,
		)

		if attempt < maxAttempts {
			base := time.Duration(attempt) * 3 * time.Second
			jitter := time.Duration(rand.Intn(2000)) * time.Millisecond
			backoff := base + jitter

			log.Printf(
				"[Cloudflare] retrying in %s...",
				backoff,
			)

			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return "", fmt.Errorf(
					"tunnel retry cancelled: %w",
					ctx.Err(),
				)
			}
		}
	}

	return "", fmt.Errorf(
		"all %d tunnel attempts failed, last error: %w",
		maxAttempts,
		lastErr,
	)
}

func (tm *Manager) GetURL() string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	return tm.PublicURL
}

func (tm *Manager) Stop() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.cancel != nil {
		tm.cancel()
	}

	tm.cmd = nil
	tm.cancel = nil
	tm.PublicURL = ""
	tm.started = false

	return nil
}