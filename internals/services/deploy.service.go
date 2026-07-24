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

func (s *DeployService) BuildDockerFile(ctx context.Context) (string, error){

	src := "/mini-paas/templates/dockerfile_node.txt"

	source, err := os.ReadFile(src)

	if err != nil {
		return "", err
	}

	fmt.Println(string(source))

	return string(source), nil

}


func (s *DeployService) CloneRepository(ctx context.Context, repoURL string, branch string) (string, error) {
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

	s.BuildDockerFile(ctx)

	return path, nil

}

