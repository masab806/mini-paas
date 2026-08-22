package dto

type LogRequest struct {
	ContainerName string `json:"container_name" required=true`
	LineNo string `json:"line_no" required=true`
}

type LogMessage struct {
	Timestamp string `json:"timestamp"`
	Log string `json:"log"`
	Source string `json:"source"`
}