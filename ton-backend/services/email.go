package services

import (
	"bytes"
	"html/template"
	"net/smtp"
	"os"
)

var isSmtpAuthenticated bool = false
var smtpAuth smtp.Auth

// Request struct
type EmailRequest struct {
	to      []string
	subject string
	body    string
}

func NewEmailRequest(to []string, subject, templateFileName string, data interface{}) (*EmailRequest, error) {
	request := &EmailRequest{
		to:      to,
		subject: subject,
		body:    "",
	}

	err := request.ParseTemplate(templateFileName, data)

	return request, err
}

func (r *EmailRequest) SendEmail() error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := os.Getenv("SMTP_HOST") + ":" + os.Getenv("SMTP_PORT")

	if !isSmtpAuthenticated {
		smtpAuth = smtp.PlainAuth("", os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_HOST"))
	}

	if err := smtp.SendMail(addr, smtpAuth, os.Getenv("SMTP_FROM"), r.to, msg); err != nil {
		return err
	}

	return nil
}

func (r *EmailRequest) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
