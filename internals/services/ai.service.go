package services

import (
	"context"
	"encoding/json"
	"fmt"
	"mini-paas/internals/repositories"

	"google.golang.org/genai"
)

type AIService struct {
	apiKey string
	client *genai.Client
	repo repositories.UserRepository
	mailService *MailService
}

type CrashAnalysisResponse struct {
	Summary   string `json:"summary"`
	Diagonsis string `json:"diagnosis"`
	Solution  string `json:"solution"`
	Severity  string `json:"severity"`
}

func NewAIService(ctx context.Context, apiKey string, repo repositories.UserRepository, mailService *MailService) (*AIService, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})

	if err != nil {
		return nil, err
	}

	return &AIService{
		apiKey: apiKey,
		client: client,
		repo: repo,
		mailService: mailService,
	}, nil
}

func (s *AIService) AnalyzeContainerLogs(ctx context.Context, logs string, email string) (*CrashAnalysisResponse, error) {
	// user, userErr := s.repo.GetByEmail(ctx, email)
	// if userErr != nil {
	// 	return nil, fmt.Errorf("failed to fetch user by email: %w", userErr)
	// }

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
			Required: []string{"summary", "diagnosis", "solution", "severity"},
		},
		SystemInstruction: genai.NewContentFromText(
			"You are an expert AI DevOps Engineer. Analyze container crash logs and identify error codes, root causes, and actionable developer fixes.",
			"",
		),
	}

	prompt := fmt.Sprintf("Analyze the following container application logs:\n\n%s", logs)

	resp, err := s.client.Models.GenerateContent(
		ctx,
		"gemini-3.6-flash",
		genai.Text(prompt),
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("gemini api call error: %w", err)
	}

	var analysis CrashAnalysisResponse
	if err := json.Unmarshal([]byte(resp.Text()), &analysis); err != nil {
		return nil, fmt.Errorf("failed to structure data in JSON: %w", err)
	}

	// if s.mailService != nil {
	// 	report := CrashReportData{
	// 		ContainerID: "app-container",
	// 		Summary:     analysis.Summary,
	// 		Diagnosis:   analysis.Diagonsis,
	// 		Solution:    analysis.Solution,
	// 		Severity:    analysis.Severity,
	// 	}

	// 	if err := s.mailService.SendCrashReport(user.Email, report); err != nil {
	// 		fmt.Printf("[AIService] Warning: failed to send crash report email: %v\n", err)
	// 	}

	return &analysis, nil
}