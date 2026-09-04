package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	baseUrl string
)

type Config struct {
	Token string `json:"token"`
}

var rootCmd = &cobra.Command{
	Use: "minipaas",
	Short: "MiniPaas CLI - Manage Deployments",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Init() {
	rootCmd.PersistentFlags().StringVar(&baseUrl, "api-url", "https://localhost:8000", "Base Url For Mini Paas Server")
}

func getConfigPath() (string, error) {
	home , err := os.UserHomeDir()
	if err != nil {
		return  "", err
	}

	return filepath.Join(home, ".minipaas.json"), nil
}

func SaveToken(token string) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	cfg := Config{Token: token}

	data, err := json.MarshalIndent(cfg, " ", "")

	if err != nil {
		return  err
	}

	return os.WriteFile(path, data, 0600)
}

func GetToken() (string, error) {
	path, err := getConfigPath()

	if err != nil {
		return "", nil
	}

	if _, err := os.Stat(path); os.IsNotExist(err){
		return "", fmt.Errorf("not logged in, Run 'minipaas login' first")
	}

	data , err := os.ReadFile(path)

	if err != nil {
		return "", nil
	}

	var cfg Config

	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", err
	}

	return cfg.Token, nil
}

func SendHttpRequest(method, endpoint string, body interface{}, requiresAuth bool) ([]byte, error) {
	var bodyReader io.Reader

	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, baseUrl+endpoint, bodyReader)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if requiresAuth {
		token, err := GetToken()
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody , err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API Error [%d]: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

