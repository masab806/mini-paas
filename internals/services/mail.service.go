package services

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

type MailService struct{}

func NewMailService() *MailService {
	return &MailService{}
}

// Struct matching your Gemini analysis payload for template binding
type CrashReportData struct {
	ContainerID string
	Summary     string
	Diagnosis   string
	Solution    string
	Severity    string
}

// SendCrashReport formats and dispatches the HTML crash report
func (s *MailService) SendCrashReport(recipientEmail string, data CrashReportData) error {
	from := os.Getenv("MY_EMAIL")
	password := os.Getenv("PASSWORD")

	if from == "" || password == "" {
		return fmt.Errorf("email credentials (MY_EMAIL / PASSWORD) are missing in environment variables")
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	to := []string{recipientEmail}

	// 1. Compile HTML Template
	tmpl, err := template.New("crash_report").Parse(crashEmailTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var bodyBytes bytes.Buffer
	if err := tmpl.Execute(&bodyBytes, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// 2. Build MIME Headers
	subject := fmt.Sprintf("[%s] Container Crash Alert: %s", data.Severity, data.ContainerID)
	
	header := make(map[string]string)
	header["From"] = fmt.Sprintf("Deployment Monitor <%s>", from)
	header["To"] = recipientEmail
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"UTF-8\""

	var message bytes.Buffer
	for k, v := range header {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.Write(bodyBytes.Bytes())

	// 3. Dispatch Email
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message.Bytes())
	if err != nil {
		return fmt.Errorf("failed to send email via SMTP: %w", err)
	}

	fmt.Printf("[MailService] Crash report successfully sent to %s\n", recipientEmail)
	return nil
}

const crashEmailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f4f5f7; margin: 0; padding: 20px; }
        .card { max-width: 600px; margin: 0 auto; background: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 6px rgba(0,0,0,0.08); border: 1px solid #e5e7eb; }
        .header { padding: 18px 24px; color: #ffffff; font-weight: 700; font-size: 16px; text-transform: uppercase; letter-spacing: 0.5px; }
        .CRITICAL { background-color: #dc2626; }
        .HIGH { background-color: #ea580c; }
        .MEDIUM { background-color: #d97706; }
        .LOW { background-color: #2563eb; }
        .content { padding: 24px; color: #1f2937; }
        .section-title { font-size: 11px; font-weight: 700; text-transform: uppercase; color: #6b7280; margin-bottom: 6px; letter-spacing: 0.5px; }
        .box { background: #f9fafb; border: 1px solid #e5e7eb; border-radius: 6px; padding: 12px 14px; margin-bottom: 18px; font-size: 14px; line-height: 1.5; color: #374151; }
        .code-box { background: #0f172a; color: #38bdf8; font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace; border-radius: 6px; padding: 14px; font-size: 13px; white-space: pre-wrap; word-break: break-word; line-height: 1.4; }
    </style>
</head>
<body>
    <div class="card">
        <div class="header {{.Severity}}">
            [{{.Severity}}] Crash Alert &bull; {{.ContainerID}}
        </div>
        <div class="content">
            <div class="section-title">Summary</div>
            <div class="box">{{.Summary}}</div>

            <div class="section-title">Diagnosis & Root Cause</div>
            <div class="box">{{.Diagnosis}}</div>

            <div class="section-title">Actionable Solution</div>
            <div class="code-box">{{.Solution}}</div>
        </div>
    </div>
</body>
</html>
`