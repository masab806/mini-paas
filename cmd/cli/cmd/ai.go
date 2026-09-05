package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type AnalyzeResponse struct {
	Severity  string `json:"severity"`
	Summary   string `json:"summary"`
	Diagnosis string `json:"diagnosis"`
	Solution  string `json:"solution"`
}

var AnalyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze container crash logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		containerName, _ := cmd.Flags().GetString("container_name")

		payload := map[string]string{
			"container_name": containerName,
		}

		response, err := SendHttpRequest("POST", "/api/ai/analyze", payload, false)
		if err != nil {
			return err
		}

		var result AnalyzeResponse
		if err := json.Unmarshal(response, &result); err != nil {
			return fmt.Errorf("failed to parse analysis response: %w", err)
		}

		cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Println()
		fmt.Println(cyan("=== CRASH ANALYSIS REPORT ==="))
		fmt.Printf("%s %s\n", bold("Severity :"), green(result.Severity))
		fmt.Printf("%s %s\n", bold("Summary  :"), result.Summary)
		fmt.Printf("%s %s\n", bold("Diagnosis:"), result.Diagnosis)
		fmt.Printf("%s %s\n", bold("Solution :"), result.Solution)
		fmt.Println()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(AnalyzeCmd)

	AnalyzeCmd.Flags().StringP("container_name", "c", "", "Container Name")
	AnalyzeCmd.MarkFlagRequired("container_name")
}