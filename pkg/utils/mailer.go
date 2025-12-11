package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

type Mailer interface {
	Send(to, subject, body string) error
}

type SMTPMailer struct {
	host     string
	port     string
	username string
	password string
	fromName string
	fromAddr string
	security string 
}

func NewSMTPMailerFromEnv() *SMTPMailer {
	security := os.Getenv("SMTP_SECURITY")
	if security == "" {
		security = "tls"
	}

	mailer := &SMTPMailer{
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
		username: os.Getenv("SMTP_USER"),
		password: os.Getenv("SMTP_PASSWORD"),
		fromName: os.Getenv("SMTP_FROM_NAME"),
		fromAddr: os.Getenv("SMTP_USER"),
		security: security,
	}

	// Debug log to check configuration
	log.Printf("SMTP Config - Host: %s, Port: %s, User: %s, Security: %s", 
		mailer.host, mailer.port, mailer.username, mailer.security)

	return mailer
}

func (m *SMTPMailer) buildMessage(toEmail, subject, body string) []byte {
	return []byte(fmt.Sprintf(
		"From: %s <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"utf-8\"\r\n"+
			"\r\n"+
			"%s\r\n",
		m.fromName, m.fromAddr, toEmail, subject, body,
	))
}

func (m *SMTPMailer) Send(toEmail, subject, body string) error {
	// Validate configuration
	if m.host == "" || m.port == "" || m.username == "" || m.password == "" {
		return fmt.Errorf("SMTP configuration is incomplete")
	}

	addr := fmt.Sprintf("%s:%s", m.host, m.port)
	msg := m.buildMessage(toEmail, subject, body)

	log.Printf("Attempting SMTP connection to %s with security: %s", addr, m.security)

	switch m.security {
	case "ssl":
		// For port 465 (SSL)
		tlsConfig := &tls.Config{
			ServerName: m.host,
			//dev mode turn it off when it comes to production
			InsecureSkipVerify: true,
		}

		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to dial SSL SMTP: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, m.host)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
		defer client.Quit()

		auth := smtp.PlainAuth("", m.username, m.password, m.host)
		if ok, _ := client.Extension("AUTH"); ok {
			if err := client.Auth(auth); err != nil {
				return fmt.Errorf("failed to auth SMTP: %w", err)
			}
		}

		if err := client.Mail(m.fromAddr); err != nil {
			return fmt.Errorf("failed to set MAIL FROM: %w", err)
		}

		if err := client.Rcpt(toEmail); err != nil {
			return fmt.Errorf("failed to set RCPT TO: %w", err)
		}

		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to start DATA: %w", err)
		}

		if _, err := w.Write(msg); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}

		return w.Close()

	case "tls":
		// For explicit TLS (usually port 587)
		tlsConfig := &tls.Config{
			ServerName: m.host,
		}

		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to dial TLS SMTP: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, m.host)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
		defer client.Quit()

		auth := smtp.PlainAuth("", m.username, m.password, m.host)
		if ok, _ := client.Extension("AUTH"); ok {
			if err := client.Auth(auth); err != nil {
				return fmt.Errorf("failed to auth SMTP: %w", err)
			}
		}

		if err := client.Mail(m.fromAddr); err != nil {
			return fmt.Errorf("failed to set MAIL FROM: %w", err)
		}

		if err := client.Rcpt(toEmail); err != nil {
			return fmt.Errorf("failed to set RCPT TO: %w", err)
		}

		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to start DATA: %w", err)
		}

		if _, err := w.Write(msg); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}

		return w.Close()

	case "starttls":
		// For STARTTLS (usually port 587)
		c, err := smtp.Dial(addr)
		if err != nil {
			return fmt.Errorf("failed to dial SMTP: %w", err)
		}
		defer c.Quit()

		tlsConfig := &tls.Config{
			ServerName: m.host,
		}

		if err := c.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}

		auth := smtp.PlainAuth("", m.username, m.password, m.host)
		if err := c.Auth(auth); err != nil {
			return fmt.Errorf("failed to auth SMTP: %w", err)
		}

		if err := c.Mail(m.fromAddr); err != nil {
			return fmt.Errorf("failed to set MAIL FROM: %w", err)
		}
		if err := c.Rcpt(toEmail); err != nil {
			return fmt.Errorf("failed to set RCPT TO: %w", err)
		}

		w, err := c.Data()
		if err != nil {
			return fmt.Errorf("failed to start DATA: %w", err)
		}

		if _, err := w.Write(msg); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}

		return w.Close()

	default:
		// Plain SMTP (not recommended for production)
		auth := smtp.PlainAuth("", m.username, m.password, m.host)
		return smtp.SendMail(addr, auth, m.fromAddr, []string{toEmail}, msg)
	}
}