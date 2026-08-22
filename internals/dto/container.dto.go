package dto

type ContainerRequest struct {
	ContainerName string `json:"container_name" binding:"required"`
}