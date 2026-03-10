package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"

	"hr-system/internal/config"
)

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
