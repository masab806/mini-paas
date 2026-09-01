package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func GenerateNginxConfig(appName string, hostPort string, domainName string) error {
	nginxConfigDir := "/etc/nginx/sites-avaliable"
	enabledDir := "/etc/nginx/sites-enabled"

	_ = os.MkdirAll(nginxConfigDir, 0755)
	_ = os.MkdirAll(enabledDir, 0755)

	configContent := fmt.Sprintf(`server {
    listen 80;
    server_name %s.%s;

    location / {
        proxy_pass http://127.0.0.1:%s;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
`, appName, domainName, hostPort)

	confPath := filepath.Join(nginxConfigDir, appName+".conf")
	symLinkPath := filepath.Join(enabledDir, appName+".conf")

	if err := os.WriteFile(confPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write nginx config: %w", err)
	}

	_ = os.Remove(symLinkPath)
	if err := os.Symlink(confPath, symLinkPath); err != nil {
		return  fmt.Errorf("failed to create nginx symlink: %w", err)
	}

	return ReloadConfig()
}

func ReloadConfig() error {
	cmd := exec.Command("nginx", "-s", "reload")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("nginx failed to load: %v, output: %s", err, string(output))
	}

	return nil
}