package utils

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"sync"

)

type TunnelManager struct {
	PublicDomain string
	cmd          *exec.Cmd
	mu            sync.RWMutex
}

var GlobalTunnel *TunnelManager

func InitTunnel(ctx context.Context, localPort string) (*TunnelManager, error) {
	tm := &TunnelManager{}

	cmd := exec.CommandContext(ctx, "cloudflared", "tunnel", "--url", fmt.Sprintf("https://localhost:%s", localPort))

	stderr, err :=  cmd.StderrPipe()

	if err != nil {
		return nil, fmt.Errorf("failed to open stderr pipe: %s", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start cmd: %s", err)
	}

	tm.cmd = cmd

	urlChan := make(chan string, 1)

	go func() {
		scanner := bufio.NewScanner(stderr);	
		re := regexp.MustCompile(`https://[a-zA-Z0-9-]+\.trycloudflare\.com`)

		for scanner.Scan(){
			line := scanner.Text()
			if match := re.FindString(line); match != "" {
				urlChan <- match
				return 
			}
		}
	}()

	select {
	case fullUrl := <- urlChan:
		domain := strings.TrimPrefix(fullUrl, "https://")
		tm.mu.Lock()
		tm.PublicDomain = domain
		tm.mu.Unlock()
		log.Printf("[Tunnel] Public HTTPS edge Active at: %s", fullUrl)

	case <- ctx.Done():
		return  nil, fmt.Errorf("tunnel setup timed out!")
	}

	GlobalTunnel = tm;
	return tm, nil

}

func (tm *TunnelManager) GetDomain() string {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	return tm.PublicDomain
}