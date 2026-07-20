package dto

type DeployRequest struct {
	RepoURL string `json:"repo_url",required`
	Branch string `json:"branch",required`
}