package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

type LogsResponse struct {
	Result string `json:"Result"` 
}

var LogCmd = &cobra.Command{
	Use:   "log",
	Short: "Logs For The Container",
	RunE: func(cmd *cobra.Command, args []string) error {
		containerName, _ := cmd.Flags().GetString("container_name")

		payload := map[string]string{
			"container_name": containerName,
		}

		response, err := SendHttpRequest("POST", "/api/logs/getLogs", payload, false)
		if err != nil {
			return err
		}

		var result LogsResponse
		if err := json.Unmarshal(response, &result); err != nil {
			return err
		}

		fmt.Println("\n==================================================================")
		fmt.Println("                    LOGS RESULT                    ")
		fmt.Println("==================================================================")
		fmt.Printf("Logs:\n%s\n", result.Result)
		fmt.Println("==================================================================")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(LogCmd)

	LogCmd.Flags().StringP("container_name", "c", "", "Container Name")
	LogCmd.MarkFlagRequired("container_name")
}