package repositories

import (
	"context"
	"mini-paas/ent"
)

type DeploymentRepository interface {
	CreateDeployment(ctx context.Context, branch string, image_tag string, repo_url string, status string, domain string, UserId int) (*ent.Deployments, error)
}
 

type deploymentRepository struct {
	client *ent.Client
}

func NewDeploymentRepository(client *ent.Client) DeploymentRepository {
	return &deploymentRepository{
		client: client,
	}
}

func (r *deploymentRepository) CreateDeployment(ctx context.Context, branch string, image_tag string, repo_url string, status string, domain string, UserId int) (*ent.Deployments, error) {
	return r.client.Deployments.
		Create().
		SetBranch(branch).
		SetImageTag(image_tag).
		SetRepoURL(repo_url).
		SetStatus(status).
		SetDomain(domain).
		SetAuthorID(UserId).
		Save(ctx)
}