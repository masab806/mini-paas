package services

import (
	"context"
	"fmt"
	"mini-paas/ent"
	"mini-paas/internals/repositories"
	"mini-paas/internals/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const DockerNetwork = "mini-paas"

type DeployService struct {
	client *ent.Client
	repo   repositories.DeploymentRepository
}

func NewDeployService(
	client *ent.Client,
	repo repositories.DeploymentRepository,
) *DeployService {
	return &DeployService{
		client: client,
		repo:   repo,
	}
}


func (s *DeployService) BuildDockerFile(
	ctx context.Context,
	projectPath string,
	imageTag string,
	framework string,
) error {

	var templateName string

	switch framework {
	case "nodejs":
		templateName = "dockerfile_node.txt"

	case "fastapi":
		templateName = "dockerfile_py.txt"

	default:
		return fmt.Errorf(
			"unsupported framework: %s",
			framework,
		)
	}

	src := filepath.Join(
		"templates",
		templateName,
	)

	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf(
			"failed to read Dockerfile template: %w",
			err,
		)
	}

	dockerfilePath := filepath.Join(
		projectPath,
		"Dockerfile",
	)

	if err := os.WriteFile(
		dockerfilePath,
		content,
		0644,
	); err != nil {
		return fmt.Errorf(
			"failed to write Dockerfile: %w",
			err,
		)
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
		return fmt.Errorf(
			"docker build failed: %w\n%s",
			err,
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}


func (s *DeployService) RunContainer(
	ctx context.Context,
	containerPort string,
	containerName string,
	imageName string,
) (string, error) {
 
	if containerPort == "" {
		return "", fmt.Errorf(
			"container port cannot be empty",
		)
	}
 
	if containerName == "" {
		return "", fmt.Errorf(
			"container name cannot be empty",
		)
	}
 
	if imageName == "" {
		return "", fmt.Errorf(
			"image name cannot be empty",
		)
	}
 
	networkCheck := exec.CommandContext(
		ctx,
		"docker",
		"network",
		"inspect",
		DockerNetwork,
	)
 
	if err := networkCheck.Run(); err != nil {
 
		createNetwork := exec.CommandContext(
			ctx,
			"docker",
			"network",
			"create",
			DockerNetwork,
		)
 
		output, err := createNetwork.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf(
				"failed to create Docker network: %w\n%s",
				err,
				strings.TrimSpace(string(output)),
			)
		}
	}
 

	cmd := exec.CommandContext(
		ctx,
		"docker",
		"run",
		"-d",
		"--name",
		containerName,
		"--network",
		DockerNetwork,
		"-e",
		fmt.Sprintf("PORT=%s", containerPort),
		imageName,
	)
 
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf(
			"docker run failed: %w\n%s",
			err,
			strings.TrimSpace(string(output)),
		)
	}
 
	containerID := strings.TrimSpace(
		string(output),
	)
 
	if containerID == "" {
		return "", fmt.Errorf(
			"docker returned an empty container ID",
		)
	}
 
	return containerID, nil
}
 


func (s *DeployService) CloneRepository(
	ctx context.Context,
	repoURL *string,
	branch *string,
	imageTag *string,
	framework *string,
	name *string,
	port *string,
	userId *int,
) (*ent.Deployments, *string, *string, error) {


	if repoURL == nil || strings.TrimSpace(*repoURL) == "" {
		return nil, nil, nil, fmt.Errorf(
			"repository URL cannot be empty",
		)
	}

	if branch == nil || strings.TrimSpace(*branch) == "" {
		return nil, nil, nil, fmt.Errorf(
			"branch cannot be empty",
		)
	}

	if imageTag == nil || strings.TrimSpace(*imageTag) == "" {
		return nil, nil, nil, fmt.Errorf(
			"image tag cannot be empty",
		)
	}

	if framework == nil || strings.TrimSpace(*framework) == "" {
		return nil, nil, nil, fmt.Errorf(
			"framework cannot be empty",
		)
	}

	if name == nil || strings.TrimSpace(*name) == "" {
		return nil, nil, nil, fmt.Errorf(
			"deployment name cannot be empty",
		)
	}

	if port == nil || strings.TrimSpace(*port) == "" {
		return nil, nil, nil, fmt.Errorf(
			"container port cannot be empty",
		)
	}

	if userId == nil {
		return nil, nil, nil, fmt.Errorf(
			"user ID cannot be nil",
		)
	}


	deploymentName := strings.TrimSpace(*name)
	containerPort := strings.TrimSpace(*port)


	baseDir := "./deployments"

	if err := os.MkdirAll(
		baseDir,
		0755,
	); err != nil {
		return nil, nil, nil, fmt.Errorf(
			"failed to create deployments directory: %w",
			err,
		)
	}

	if _, err := exec.LookPath("git"); err != nil {
		return nil, nil, nil, fmt.Errorf(
			"git not found: %w",
			err,
		)
	}


	deploymentID := uuid.NewString()

	repoPath := filepath.Join(
		baseDir,
		deploymentID,
	)

	cloneCmd := exec.CommandContext(
		ctx,
		"git",
		"clone",
		"--branch",
		*branch,
		"--single-branch",
		*repoURL,
		repoPath,
	)

	output, err := cloneCmd.CombinedOutput()
	if err != nil {
		return nil, nil, nil, fmt.Errorf(
			"git clone failed: %w\n%s",
			err,
			strings.TrimSpace(string(output)),
		)
	}


	if err := s.BuildDockerFile(
		ctx,
		repoPath,
		*imageTag,
		*framework,
	); err != nil {

		_ = os.RemoveAll(repoPath)

		return nil, nil, nil, fmt.Errorf(
			"failed to build Docker image: %w",
			err,
		)
	}

	containerID, err := s.RunContainer(
		ctx,
		containerPort,
		deploymentName,
		*imageTag,
	)

	if err != nil {
		_ = os.RemoveAll(repoPath)

		return nil, nil, nil, fmt.Errorf(
			"failed to run container: %w",
			err,
		)
	}


	dockerStatus, err := utils.GetContainerStatus(
		ctx,
		containerID,
	)

	if err != nil {

		_ = stopAndRemoveContainer(
			ctx,
			deploymentName,
		)

		_ = os.RemoveAll(repoPath)

		return nil, nil, nil, fmt.Errorf(
			"failed to get container status: %w",
			err,
		)
	}


	route := utils.NginxRoute{
		Domain:        deploymentName,
		ContainerName: deploymentName,
		Port:          containerPort,
	}

	if err := configureNginxRoute(
		route,
	); err != nil {

		_ = stopAndRemoveContainer(
			ctx,
			deploymentName,
		)

		_ = os.RemoveAll(repoPath)

		return nil, nil, nil, fmt.Errorf(
			"failed to configure nginx: %w",
			err,
		)
	}



	publicBaseURL, err := utils.GlobalTunnel.StartCloudflareTunnelWithRetry(
		context.Background(),
		5,              
		30*time.Second, 
	)

	if err != nil {

		_ = stopAndRemoveContainer(
			ctx,
			deploymentName,
		)

		_ = os.RemoveAll(repoPath)

		return nil, nil, nil, fmt.Errorf(
			"failed to start Cloudflare tunnel: %w",
			err,
		)
	}


	publicURL := fmt.Sprintf(
		"%s/%s/",
		strings.TrimRight(
			publicBaseURL,
			"/",
		),
		deploymentName,
	)

	deployment, err := s.repo.CreateDeployment(
		ctx,
		*branch,
		*imageTag,
		*repoURL,
		dockerStatus,
		publicURL,
		*userId,
	)

	if err != nil {

		_ = stopAndRemoveContainer(
			ctx,
			deploymentName,
		)

		_ = os.RemoveAll(repoPath)

		return nil, nil, nil, fmt.Errorf(
			"failed to save deployment: %w",
			err,
		)
	}

	return deployment, &repoPath, &publicURL, nil
}

func configureNginxRoute(
	route utils.NginxRoute,
) error {

	if strings.TrimSpace(route.Domain) == "" {
		return fmt.Errorf(
			"nginx domain cannot be empty",
		)
	}

	if strings.TrimSpace(route.ContainerName) == "" {
		return fmt.Errorf(
			"nginx container name cannot be empty",
		)
	}

	if strings.TrimSpace(route.Port) == "" {
		return fmt.Errorf(
			"nginx port cannot be empty",
		)
	}


	return utils.GenerateNginxConfig(
		[]utils.NginxRoute{
			route,
		},
	)
}


func stopAndRemoveContainer(
	ctx context.Context,
	containerName string,
) error {

	stopCmd := exec.CommandContext(
		ctx,
		"docker",
		"rm",
		"-f",
		containerName,
	)

	output, err := stopCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"failed to remove container: %w\n%s",
			err,
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}
