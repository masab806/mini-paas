package services

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
)

type LogService struct{}

func NewLogService() *LogService {
	return &LogService{}
}

var npmBoilerplate = regexp.MustCompile(`(?m)^>\s*.*$\n?`)

func SimpleSanitize(raw string) string {
	clean := npmBoilerplate.ReplaceAllString(raw, "")

	lines := strings.Split(clean, "\n")
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return strings.Join(result, "\n")
}

func (s *LogService) GetContainerLogs(ctx context.Context, containerName string, lineNo string) (string, error) {
	cmd := exec.CommandContext(
		ctx,
		"docker",
		"logs",
		"--tail",
		lineNo,
		containerName,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", err
	}

	finalResult := SimpleSanitize(string(output))

	return finalResult, nil
}

func (s *LogService) StreamContainerLogs(ctx context.Context, containerName string, lineNo string) (io.ReadCloser, error) {
	cmd := exec.CommandContext(
		ctx,
		"docker",
		"logs",
		"-f",
		"--tail",
		lineNo,
		containerName,
	)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return nil, fmt.Errorf("failed to create a pipeline: %w", err)
	}

	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start log command: %w", err)
	}

	return stdout, nil
}