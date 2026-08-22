package services

import (
	"context"
	"fmt"
	"os/exec"
)

type ContainerService struct{}

func NewContainerService() *ContainerService {
	return &ContainerService{}
}

func (s *ContainerService) RunContainer(ctx context.Context, containerName string) (string, error){
	cmd := exec.CommandContext(
		ctx,
		"docker",
		"start",
		containerName,
	)

	output, err := cmd.CombinedOutput()

	fmt.Printf(string(output), err)

	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (s *ContainerService) StopContainer(ctx context.Context, containerName string) (string, error) {
	cmd := exec.CommandContext(
		ctx,
		"docker",
		"stop",
		containerName,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (s *ContainerService) DeleteContainer(ctx context.Context, containerName string) (string, error){
	cmd := exec.CommandContext(
		ctx,
		"docker",
		"rm",
		containerName,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return  "", err
	}

	return string(output), nil
}