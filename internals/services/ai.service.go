package services

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/genai"
)

type AIService struct {
	apiKey string
	client *genai.Client
}

type CrashAnalysisResponse struct {
	Summary   string `json:"summary"`
	Diagonsis string `json:"diagnosis"`
	Solution  string `json:"solution"`
	Severity  string `json:"severity"`
}

func NewAIService(ctx context.Context, apiKey string) (*AIService, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})

	if err != nil {
		return nil, err
	}

	return &AIService{
		apiKey: apiKey,
		client: client,
	}, nil
}

func (s *AIService) AnalyzeContainerLogs(ctx context.Context, logs string) (*CrashAnalysisResponse, error) {
	var promptTemplate = `You are an expert AI DevOps Engineer and System Reliability Specialist. Analyze the provided container application crash logs or deployment stack traces and output a JSON object strictly matching the following schema.

### JSON Output Schema:
{
  "summary": "Brief 1-2 sentence overview of what went wrong.",
  "diagnosis": "Detailed root-cause analysis identifying specific error codes, failing components, memory/CPU issues, missing environment variables, or code faults.",
  "solution": "Step-by-step actionable guide or code fix to resolve the issue permanently.",
  "severity": "One of: LOW, MEDIUM, HIGH, CRITICAL"
}

### Guidelines:
1. Ensure the JSON is valid, well-formed, and strictly adheres to the exact key names shown above.
2. Do not include markdown code block syntax (like ` + "```json" + `) in your response unless raw JSON output is required.
3. Keep the "severity" strictly capped to one of the four allowed values: "LOW", "MEDIUM", "HIGH", or "CRITICAL".
4. Focus on practical, developer-centric solutions (e.g., exact command line fixes, environment configuration updates, or code adjustments).

### Logs:
%s`

	prompt := fmt.Sprintf(promptTemplate, logs)

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"summary": {
					Type:        genai.TypeString,
					Description: "A concise 1-2 sentence summary of the container crash event.",
				},
				"diagnosis": {
					Type:        genai.TypeString,
					Description: "Detailed technical analysis explaining the root cause of the error.",
				},
				"solution": {
					Type:        genai.TypeString,
					Description: "Actionable steps or code fixes to resolve the issue.",
				},
				"severity": {
					Type:        genai.TypeString,
					Description: "The severity level of the crash event.",
					Enum:        []string{"LOW", "MEDIUM", "HIGH", "CRITICAL"},
				},
			},
			Required: []string{"summary", "diagnosis", "solution"},
		},
		SystemInstruction: genai.NewContentFromText(
			"Analyze the container logs and return a concise summary, detailed diagnosis, and actionable solution.",
			genai.RoleUser,
		),
	}

	resp, err := s.client.Models.GenerateContent(
		ctx,
		"gemini-3.6-flash",
		genai.Text(prompt),
		config,
	)

	if err != nil {
		return nil, fmt.Errorf("gemini api call error: %s", err)
	}

	var analysis CrashAnalysisResponse

	if err := json.Unmarshal([]byte(resp.Text()), &analysis); err != nil {
		return nil, fmt.Errorf("failed to structure data in JSON: %s", err)
	}

	return &analysis, nil
}
