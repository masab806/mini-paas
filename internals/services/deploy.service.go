package services

import (
	"context"
	"errors"
	"fmt"
	"mini-paas/ent"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
)

type DeployService struct {
	client *ent.Client
}

func NewDeployService(client *ent.Client) *DeployService {
	return &DeployService{
		client: client,
	}
}

func (s *DeployService) BuildDockerFile(ctx context.Context, projectPath string, imageTag string, framework string) error {

	var templateName string

	switch framework {
	case "nodejs":
		templateName = "dockerfile_node.txt"

	case "fastapi":
		templateName = "dockerfile_py.txt"

	default:
		return fmt.Errorf("Unsupported Frameworks: %s", framework)
	}

	src := filepath.Join("/mini-paas/templates", templateName)

	content, err := os.ReadFile(src)

	if err != nil {
		return err
	}

	dockerfilePath := filepath.Join(projectPath, "Dockerfile")

	err = os.WriteFile(dockerfilePath, content, 0644)

	if err != nil {
		return err
	}

	cmd := exec.CommandContext(
		ctx,
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

func (s *DeployService) RunContainer(ctx context.Context, port string, containerName string, imageName string) error {
	cmd := exec.CommandContext(
		ctx,
		"docker",
		"run",
		"-d",
		"--name", containerName, 
		"-p", port, 
		imageName,
	)



	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	fmt.Println(string(output))

	return nil

}

func (s *DeployService) CloneRepository(ctx context.Context, repoURL string, branch string, imageTag string, framework string, name string, port string) (string, error) {
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

	err = s.BuildDockerFile(ctx, path, imageTag, framework)

	containerErr := s.RunContainer(ctx, port, name, imageTag)

	if containerErr != nil {
		return "", containerErr
	}

	if err != nil {
		return "", err
	}

	return path, nil

}
