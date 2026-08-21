package dto

type LogRequest struct {
	ContainerName string `json:"container_name" required=true`
	LineNo string `json:"line_no" required=true`
}