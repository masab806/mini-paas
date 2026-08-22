package services

import (
	"context"
	"errors"
	"fmt"
	"mini-paas/ent"
	"mini-paas/internals/repositories"
	"mini-paas/internals/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type DeployService struct {
	client *ent.Client
	repo   repositories.DeploymentRepository
}

func NewDeployService(client *ent.Client, repo repositories.DeploymentRepository) *DeployService {
	return &DeployService{
		client: client,
		repo:   repo,
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

	return nil

}

func (s *DeployService) RunContainer(ctx context.Context, port string, containerName string, imageName string) (string, error) {
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
		return "", err
	}

	fmt.Println(string(output))

	containerID := strings.TrimSpace(string(output))

	return containerID, nil

}

func (s *DeployService) CloneRepository(ctx context.Context, repoURL string, branch string, imageTag string, framework string, name string, port string, userId int) (*ent.Deployments, string, error) {
	baseDir := "/mini-pass/deployments/"

	err := os.MkdirAll(baseDir, 0755)

	if err != nil {
		return nil, "", errors.New("Invalid Path")
	}

	id := uuid.NewString()

	path := filepath.Join(baseDir, id)

	cmd := exec.Command("git", "clone", repoURL, path)

	gitPath, gitErr := exec.LookPath("git")

	if err != nil {
		return nil, "", gitErr
	}

	fmt.Println(gitPath)

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("git clone failed: %v\n", err)
		fmt.Printf("output:\n%s\n", output)
		return nil, "", err
	}

	err = s.BuildDockerFile(ctx, path, imageTag, framework)

	containerId, containerErr := s.RunContainer(ctx, port, name, imageTag)

	dockerStatus, statusErr := utils.GetContainerStatus(ctx, containerId)

	if statusErr != nil {
		return nil, "", statusErr
	}

	deployment, err := s.repo.CreateDeployment(ctx, branch, imageTag, repoURL, dockerStatus, userId)

	if containerErr != nil {
		return nil, "", containerErr
	}

	if err != nil {
		return nil, "", err
	}

	return deployment, path, nil

}
