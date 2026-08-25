package dto

type AnalyzeLogResponse struct {
	ContainerName string `json:"container_name"`
	LineNo string `json:"line_no"`
}