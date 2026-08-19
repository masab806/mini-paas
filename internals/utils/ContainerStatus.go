package utils

import (
	"context"
	"os/exec"
	"strings"
)

func GetContainerStatus(ctx context.Context, containerID string) (string, error) {
	cmd := exec.CommandContext(ctx, "docker", "inspect", "--format", "{{.State.Status}}", containerID)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}