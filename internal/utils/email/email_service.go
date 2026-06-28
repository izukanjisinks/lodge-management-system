package email

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"

	"lodge-system/internal/config"
)

// Attachment represents a file to attach to an email.
type Attachment struct {
	Filename    string
	ContentType string // e.g. "application/pdf"
	Data        []byte
}

type EmailService struct {
	config *config.EmailConfig
}

func NewEmailService(cfg *config.EmailConfig) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

// SendEmail sends an email with HTML body
func (s *EmailService) SendEmail(to []string, subject, htmlBody string) error {
	log.Printf("[EMAIL] Sending email to %v with subject: %s", to, subject)

	// Build email headers
	from := fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = strings.Join(to, ", ")
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"
	headers["Date"] = time.Now().Format(time.RFC1123Z)

	// Build message
	var message strings.Builder
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.WriteString(htmlBody)

	// Get authentication
	auth := s.getAuth()

	// Connect and send
	addr := s.config.Host + ":" + s.config.Port
	if s.config.UseTLS {
		return s.sendWithTLS(addr, auth, s.config.FromEmail, to, []byte(message.String()))
	}

	return smtp.SendMail(addr, auth, s.config.FromEmail, to, []byte(message.String()))
}

// SendEmailWithAttachment sends an HTML email with one or more file attachments
// using a multipart/mixed MIME message built entirely from the standard library.
func (s *EmailService) SendEmailWithAttachment(to []string, subject, htmlBody string, attachments ...Attachment) error {
	log.Printf("[EMAIL] Sending email with %d attachment(s) to %v with subject: %s", len(attachments), to, subject)

	if len(attachments) == 0 {
		return s.SendEmail(to, subject, htmlBody)
	}

	var message strings.Builder
	writer := multipart.NewWriter(&message)

	from := fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	message.WriteString(fmt.Sprintf("From: %s\r\n", from))
	message.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ", ")))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	message.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", writer.Boundary()))
	message.WriteString("\r\n")

	// HTML body part
	htmlHeader := textproto.MIMEHeader{}
	htmlHeader.Set("Content-Type", "text/html; charset=UTF-8")
	htmlHeader.Set("Content-Transfer-Encoding", "quoted-printable")
	htmlPart, err := writer.CreatePart(htmlHeader)
	if err != nil {
		return fmt.Errorf("failed to create HTML part: %w", err)
	}
	if _, err := htmlPart.Write([]byte(htmlBody)); err != nil {
		return fmt.Errorf("failed to write HTML part: %w", err)
	}

	// Attachment parts
	for _, att := range attachments {
		ct := att.ContentType
		if ct == "" {
			ct = "application/octet-stream"
		}
		attHeader := textproto.MIMEHeader{}
		attHeader.Set("Content-Type", fmt.Sprintf("%s; name=%q", ct, att.Filename))
		attHeader.Set("Content-Transfer-Encoding", "base64")
		attHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", att.Filename))
		attPart, err := writer.CreatePart(attHeader)
		if err != nil {
			return fmt.Errorf("failed to create attachment part: %w", err)
		}
		// base64-encode in 76-char lines per RFC 2045
		encoded := base64.StdEncoding.EncodeToString(att.Data)
		for i := 0; i < len(encoded); i += 76 {
			end := i + 76
			if end > len(encoded) {
				end = len(encoded)
			}
			if _, err := attPart.Write([]byte(encoded[i:end] + "\r\n")); err != nil {
				return fmt.Errorf("failed to write attachment part: %w", err)
			}
		}
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to finalize multipart message: %w", err)
	}

	auth := s.getAuth()
	addr := s.config.Host + ":" + s.config.Port
	if s.config.UseTLS {
		return s.sendWithTLS(addr, auth, s.config.FromEmail, to, []byte(message.String()))
	}
	return smtp.SendMail(addr, auth, s.config.FromEmail, to, []byte(message.String()))
}

// getAuth returns the appropriate SMTP auth based on config
func (s *EmailService) getAuth() smtp.Auth {
	if s.config.AuthMethod == "LOGIN" {
		return &loginAuth{
			username: s.config.Username,
			password: s.config.Password,
		}
	}
	// Default to PLAIN auth
	return smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
}

// sendWithTLS sends email with STARTTLS
func (s *EmailService) sendWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	if err := client.StartTLS(&tls.Config{ServerName: s.config.Host}); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	if err := client.Mail(from); err != nil {
		return fmt.Errorf("MAIL command failed: %w", err)
	}

	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("RCPT command failed for %s: %w", recipient, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA command failed: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	log.Printf("[EMAIL] Email sent successfully to %v", to)
	return client.Quit()
}

// loginAuth implements LOGIN authentication (for legacy servers)
type loginAuth struct {
	username, password string
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		prompt := string(fromServer)
		switch {
		case strings.Contains(strings.ToLower(prompt), "username"):
			return []byte(a.username), nil
		case strings.Contains(strings.ToLower(prompt), "password"):
			return []byte(a.password), nil
		default:
			return nil, fmt.Errorf("unexpected server prompt: %s", prompt)
		}
	}
	return nil, nil
}

// TestConnection tests the SMTP connection
func (s *EmailService) TestConnection() error {
	addr := s.config.Host + ":" + s.config.Port
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	if s.config.UseTLS {
		if err := client.StartTLS(&tls.Config{ServerName: s.config.Host}); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	auth := s.getAuth()
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	log.Println("[EMAIL] SMTP connection test successful")
	return client.Quit()
}
