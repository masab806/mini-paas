package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
)

type DeployService struct{}

func NewDeployService() *DeployService {
	return &DeployService{}
}

func (s *DeployService) BuildDockerFile(ctx context.Context, projectPath string, imageTag string) error {

	src := "/mini-paas/templates/dockerfile_node.txt"

	content, err := os.ReadFile(src)

	if err != nil {
		return err
	}

	dockerfilePath := filepath.Join(projectPath, "Dockerfile")

	err = os.WriteFile(dockerfilePath, content, 0644)

	if err != nil {
		return err
	}

	cmd := exec.Command(
		"docker",
		"build",
		"-t",
		imageTag,
		".",
	)

	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("docker build failed: %w\n%s", err, output)
	}

	fmt.Println(string(output))

	return nil

}

func (s *DeployService) CloneRepository(ctx context.Context, repoURL string, branch string, imageTag string) (string, error) {
	baseDir := "/mini-pass/deployments/"

	err := os.MkdirAll(baseDir, 0755)

	if err != nil {
		return "", errors.New("Invalid Path")
	}

	id := uuid.NewString()

	path := filepath.Join(baseDir, id)

	cmd := exec.Command("git", "clone", repoURL, path)

	gitPath, gitErr := exec.LookPath("git")

	if err != nil {
		return "", gitErr
	}

	fmt.Println(gitPath)

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("git clone failed: %v\n", err)
		fmt.Printf("output:\n%s\n", output)
		return "", err
	}

	err = s.BuildDockerFile(ctx, path, imageTag)

	if err != nil {
		return "", err
	}

	return path, nil

}
