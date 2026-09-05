package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	NginxContainerName = "mini-paas-nginx"
	NginxConfigFile    = "./nginx/conf.d/mini-paas.conf"
)

type NginxRoute struct {
	Domain        string
	ContainerName string
	Port          string
}

func sanitizeVarName(s string) string {
	r := strings.NewReplacer("-", "_", ".", "_", "/", "_")
	return r.Replace(s)
}

func SetupNginx() error {
	if err := os.MkdirAll("./nginx/conf.d", 0755); err != nil {
		return fmt.Errorf(
			"failed to create nginx config directory: %w",
			err,
		)
	}

	if _, err := os.Stat(NginxConfigFile); os.IsNotExist(err) {
		defaultConfig := `server {
    listen 80;
    server_name _;

    resolver 127.0.0.11 valid=10s ipv6=off;

    location / {
        return 404 "Mini-PaaS Router: Route not found";
    }
}
`

		if err := os.WriteFile(
			NginxConfigFile,
			[]byte(defaultConfig),
			0644,
		); err != nil {
			return fmt.Errorf(
				"failed to create nginx config: %w",
				err,
			)
		}
	}

	check := exec.Command(
		"docker",
		"inspect",
		NginxContainerName,
	)

	if err := check.Run(); err != nil {
		cmd := exec.Command(
			"docker",
			"compose",
			"up",
			"-d",
			"nginx",
		)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf(
				"failed to start nginx container: %s",
				strings.TrimSpace(string(output)),
			)
		}
	} else {
		cmd := exec.Command(
			"docker",
			"inspect",
			"-f",
			"{{.State.Running}}",
			NginxContainerName,
		)

		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf(
				"failed to check nginx status: %w",
				err,
			)
		}

		if strings.TrimSpace(string(output)) != "true" {
			cmd := exec.Command(
				"docker",
				"start",
				NginxContainerName,
			)

			output, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf(
					"failed to start nginx: %s",
					strings.TrimSpace(string(output)),
				)
			}
		}
	}

	return ReloadNginx()
}

func GenerateNginxConfig(routes []NginxRoute) error {
	var config strings.Builder

	config.WriteString(`server {
    listen 80;
    server_name _;

    resolver 127.0.0.11 valid=10s ipv6=off;

    location / {
        return 404 "Mini-PaaS Router: Route not found";
    }

`)

	for _, route := range routes {
		if strings.TrimSpace(route.Domain) == "" {
			return fmt.Errorf("nginx domain cannot be empty")
		}

		if strings.TrimSpace(route.ContainerName) == "" {
			return fmt.Errorf("nginx container name cannot be empty")
		}

		if strings.TrimSpace(route.Port) == "" {
			return fmt.Errorf("nginx port cannot be empty")
		}

		domain := strings.Trim(route.Domain, "/")
		containerName := strings.TrimSpace(route.ContainerName)
		port := strings.TrimSpace(route.Port)
		varName := sanitizeVarName(domain)

		config.WriteString(
			fmt.Sprintf(`    location /%s/ {
        set $upstream_%s http://%s:%s/;
        proxy_pass $upstream_%s;

        proxy_http_version 1.1;

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        proxy_read_timeout 300s;
        proxy_send_timeout 300s;
    }

`, domain, varName, containerName, port, varName),
		)
	}

	config.WriteString(`}
`)

	tempFile := NginxConfigFile + ".tmp"

	if err := os.WriteFile(
		tempFile,
		[]byte(config.String()),
		0644,
	); err != nil {
		return fmt.Errorf(
			"failed to write temporary nginx config: %w",
			err,
		)
	}

	if err := os.Rename(
		tempFile,
		NginxConfigFile,
	); err != nil {
		_ = os.Remove(NginxConfigFile)
		if err := os.Rename(tempFile, NginxConfigFile); err != nil {
			return fmt.Errorf(
				"failed to replace nginx config: %w",
				err,
			)
		}
	}

	return ReloadNginx()
}

func ReloadNginx() error {
	test := exec.Command(
		"docker",
		"exec",
		NginxContainerName,
		"nginx",
		"-t",
	)

	output, err := test.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"nginx configuration test failed: %s",
			strings.TrimSpace(string(output)),
		)
	}

	reload := exec.Command(
		"docker",
		"exec",
		NginxContainerName,
		"nginx",
		"-s",
		"reload",
	)

	output, err = reload.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"failed to reload nginx: %s",
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}