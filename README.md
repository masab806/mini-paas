# 🚀 Mini-PaaS

A lightweight, high-performance **Platform-as-a-Service (PaaS)** engine built with **Go**, **Gin Gonic**, **Docker Engine API**, **Nginx**, **Cloudflare**, and **Cobra CLI**. Mini-PaaS automates repository cloning, container builds, zero-downtime path-based reverse proxy routing, real-time log streaming, and AI-powered container crash diagnostics.

---

## 🌟 Key Features

- **Automated Deployment Pipeline** — clones source code, builds Docker images, and spins up containers.
- **Zero-Downtime Dynamic Routing** — embedded Nginx reverse proxy with Docker's internal DNS resolver.
- **Cloudflare Public Exposure** — secure HTTPS access via Cloudflare Tunnel.
- **Resilient Proxy Config** — atomic config writes, safe across Windows/Linux.
- **Log Management** — tail or stream stdout/stderr from any app container.
- **AI Crash Analysis** — parses failure logs into severity ratings, root causes, and fixes.
- **Developer CLI** — Cobra + Viper based client for deployments, logs, and diagnostics.

---

## 🏗️ System Overview

```
                  +-------------------+
                  |   Mini-PaaS CLI   |
                  +---------+---------+
                            |
                            v
                  +-------------------+
                  |   Gin API Engine  |
                  +----+--------+-----+
                       |        |
     +-----------------+        +-----------------+
     |                                            |
     v                                            v
+-------------------+                  +-------------------+
|   Docker Engine    |                  | Nginx + Cloudflare|
| - Build & Deploy   |                  |  Reverse Proxy    |
| - Fetch Logs        |                  |                    |
+-------------------+                  +-------------------+
```

---

## 🛠️ Prerequisites

- Go `v1.20+`
- Docker Engine & Docker Compose `v20.10+`
- Git (available on PATH)

---

## 📦 Setup

Everything runs through three steps: start the proxy, start the API, then build the CLI.

```bash
# 1. Start the Nginx reverse proxy
docker compose up -d nginx

# 2. Start the API engine
go mod tidy && go run main.go

# 3. Build the CLI
go build -o mini-paas main.go        # Linux/macOS
sudo mv mini-paas /usr/local/bin/    # optional: add to PATH

# Windows (PowerShell)
go build -o mini-paas.exe main.go
```

The `docker-compose.yml` defines the Nginx gateway on an isolated `mini-paas` network — see [`nginx/`](./nginx) for config details.

---

## 💻 Usage

| Command | Description |
|---|---|
| `mini-paas deploy -r <repo> -b <branch> -t <tag> -f <framework> -p <port>` | Clone, build, deploy, and route a new app |
| `mini-paas log -c <container>` | Fetch or stream logs from a running container |
| `mini-paas diagnose -c <container>` | Run AI crash analysis on container failure logs |

Run `mini-paas --help` or `mini-paas <command> --help` for full flag reference.

---
