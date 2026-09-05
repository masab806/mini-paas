package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

type DeployResponse struct {
	Message       string      `json:"message"`
	Path          string      `json:"path"`
	PublicURL     string      `json:"public_url"`
	ContainerName string      `json:"container_name"`
	Details       interface{} `json:"details"`
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy application to web",
	RunE: func(cmd *cobra.Command, args []string) error {
		repoURL, _ := cmd.Flags().GetString("repo")
		branch, _ := cmd.Flags().GetString("branch")
		imageTag, _ := cmd.Flags().GetString("tag")
		framework, _ := cmd.Flags().GetString("framework")
		port, _ := cmd.Flags().GetString("port")

		payload := map[string]string{
			"repo_url":  repoURL,
			"branch":    branch,
			"image_tag": imageTag,
			"framework": framework,
			"port":      port,
		}

		fmt.Println("Starting Deployment Pipeline...")
		fmt.Printf("   ├─ Repository : %s [%s]\n", repoURL, branch)
		fmt.Printf("   ├─ Framework  : %s\n", framework)
		fmt.Printf("   └─ Port       : %s\n\n", port)

		fmt.Print("[1/3] Cloning repository & building image... \n")

		response, err := SendHttpRequest("POST", "/api/deploy/", payload, true)
		if err != nil {
			return err
		}
		fmt.Println("Done... \n")

		fmt.Println("[2/3] Spinning up container instance... Done! \n")
		fmt.Println("[3/3] Registering route with reverse proxy... Done! \n")

		var result DeployResponse
		if err := json.Unmarshal(response, &result); err != nil {
			return err
		}

		fmt.Println("\n==================================================================")
		fmt.Println("                    DEPLOYMENT SUCCESSFUL                     ")
		fmt.Println("==================================================================")
		fmt.Printf("Container Name: %s\n", result.ContainerName)
		fmt.Printf("🌐 Public URL : %s\n", result.PublicURL)
		fmt.Println("==================================================================")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringP("repo", "r", "", "Git Repository URL")
	deployCmd.Flags().StringP("branch", "b", "main", "Repository Branch")
	deployCmd.Flags().StringP("tag", "t", "latest", "Docker Image Tag")
	deployCmd.Flags().StringP("framework", "f", "", "Application Framework")
	deployCmd.Flags().StringP("port", "p", "", "Application Exposed Port")

	deployCmd.MarkFlagRequired("repo")
	deployCmd.MarkFlagRequired("framework")
	deployCmd.MarkFlagRequired("port")
}
