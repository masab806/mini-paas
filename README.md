# 🚀 SentinelPaaS

**Self-Orchastrated Mini-Paas** is a self-hosted, event-driven Mini-PaaS (Platform as a Service) built in Go. It automates container lifecycle management using the Docker Engine API, dynamically routes domain traffic via an Nginx reverse proxy, and features an automated **AIOps diagnostic engine**.

When a container crashes, SentinelPaaS intercepts the event, extracts trailing runtime logs, performs root-cause analysis using the Gemini API with structured JSON output, and dispatches HTML post-mortem reports directly to the application owner.

---

## ✨ Features

- **Container Lifecycle Management:** Dynamically provision, isolate, and manage user container runtimes via the Docker SDK.
- **Dynamic Reverse Proxying:** Route external HTTP/HTTPS traffic to running containers with Nginx.
- **Event-Driven Crash Monitoring:** Asynchronous background worker (`docker.Client.Events`) listening for non-zero container exit (`die`) events in real-time.
- **Automated AIOps Diagnostics:** Uses Google Gen AI SDK to analyze crash logs and output structured diagnostic payloads (`summary`, `diagnosis`, `solution`, `severity`).
- **Multi-Tenant Email Alerts:** Automatically resolves container ownership via metadata labels and sends color-coded HTML diagnostic reports to specific users.
- **Crash Loop Debouncing:** Thread-safe in-memory rate limiting to prevent spam and save API quota during continuous container restarts (`CrashLoopBackOff`).

---

## 🏗️ Architecture Overview

```text
  [ User App Container ]
            │ (Crashes / Non-Zero Exit)
            ▼
┌──────────────────────────────────────────────────────────┐
│                   Docker Engine API                      │
└───────────────────────────┬──────────────────────────────┘
                            │ (Die Event Stream)
                            ▼
┌──────────────────────────────────────────────────────────┐
│            SentinelPaaS Event Monitor Worker             │
│  - Extracts Container ID & User Metadata                 │
│  - Checks Thread-Safe Debouncer Lock                      │
│  - Fetches Trailing Stderr/Stdout Logs                   │
└───────────────────────────┬──────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────┐
│                   Gemini AI Engine                       │
│  - Enforces Structured JSON Output Schema                │
│  - Generates Diagnosis, Root Cause, and Code Solutions   │
└───────────────────────────┬──────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────┐
│                  Transactional Mailer                    │
│  - Constructs Styled HTML Email Payload                  │
│  - Dispatches Incident Report to Container Owner         │
└──────────────────────────────────────────────────────────┘
