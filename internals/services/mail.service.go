package services

import (
	"fmt"
	"net/smtp"
	"os"
)

type MailService struct{}

func NewMailService() *MailService {
	return &MailService{}
}

func (s *MailService) SendMail(receipientEmail string, messageText string) {
	from := os.Getenv("MY_EMAIL")
	password := os.Getenv("PASSWORD")

	to := []string{receipientEmail}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte(messageText)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort,auth,from,to,message)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Email Sent!")
}