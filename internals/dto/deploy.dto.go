package dto

type DeployRequest struct {
	RepoURL string `json:"repo_url",required`
	Branch string `json:"branch",required`
	ImageTag string `json:"image_tag",required`
}