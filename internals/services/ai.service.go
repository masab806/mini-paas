package services

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/genai"
)

type CrashAnalysisResponse struct {
	Summary      string `json:"summary"`
	RootCause    string `json:"root_cause"`
	SuggestedFix string `json:"suggested_fix"`
	Severity     string `json:"severity"`
}

type AIService struct {
	client *genai.Client
}

func NewAIService(ctx context.Context, apiKey string) (*AIService, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize genai client: %w", err)
	}
	return &AIService{client: client}, nil
}

func (s *AIService) AnalyzeLogs(ctx context.Context, logs string) (*CrashAnalysisResponse, error) {
	prompt := fmt.Sprintf(`You are an expert DevOps engineer.
Analyze these application container logs from a crashed service and return a diagnosis:

%s`, logs)

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"summary":       {Type: genai.TypeString},
				"root_cause":    {Type: genai.TypeString},
				"suggested_fix": {Type: genai.TypeString},
				"severity":      {Type: genai.TypeString},
			},
			Required: []string{"summary", "root_cause", "suggested_fix", "severity"},
		},
	}

	result, err := s.client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("gemini generation failed: %w", err)
	}

	var analysis CrashAnalysisResponse
	if err := json.Unmarshal([]byte(result.Text()), &analysis); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &analysis, nil
}