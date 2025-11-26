package email

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

// Config holds email service configuration
type Config struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromAddress  string
	FromName     string
	TemplateDir  string
}

// Service handles email operations
type Service struct {
	config    *Config
	dialer    *gomail.Dialer
	templates map[string]*template.Template
}

// NewService creates a new email service
func NewService(config *Config) (*Service, error) {
	dialer := gomail.NewDialer(
		config.SMTPHost,
		config.SMTPPort,
		config.SMTPUser,
		config.SMTPPassword,
	)

	service := &Service{
		config:    config,
		dialer:    dialer,
		templates: make(map[string]*template.Template),
	}

	// Load email templates
	if err := service.loadTemplates(); err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	return service, nil
}

// loadTemplates loads all email templates
func (s *Service) loadTemplates() error {
	templateNames := []string{
		"verification",
		"password_reset",
		"welcome",
	}

	for _, name := range templateNames {
		templatePath := filepath.Join(s.config.TemplateDir, name+".html")
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			// If template file doesn't exist, use inline template
			tmpl, err = s.getInlineTemplate(name)
			if err != nil {
				return fmt.Errorf("failed to load template %s: %w", name, err)
			}
		}
		s.templates[name] = tmpl
	}

	return nil
}

// getInlineTemplate returns an inline HTML template
func (s *Service) getInlineTemplate(name string) (*template.Template, error) {
	var templateHTML string

	switch name {
	case "verification":
		templateHTML = verificationTemplate
	case "password_reset":
		templateHTML = passwordResetTemplate
	case "welcome":
		templateHTML = welcomeTemplate
	default:
		return nil, fmt.Errorf("unknown template: %s", name)
	}

	return template.New(name).Parse(templateHTML)
}

// SendVerificationEmail sends an email verification email
func (s *Service) SendVerificationEmail(to, displayName, verificationURL string) error {
	data := map[string]interface{}{
		"DisplayName":     displayName,
		"VerificationURL": verificationURL,
	}

	subject := "Verify your NexusFlow account"
	return s.sendEmail(to, subject, "verification", data)
}

// SendPasswordResetEmail sends a password reset email
func (s *Service) SendPasswordResetEmail(to, displayName, resetURL string) error {
	data := map[string]interface{}{
		"DisplayName": displayName,
		"ResetURL":    resetURL,
	}

	subject := "Reset your NexusFlow password"
	return s.sendEmail(to, subject, "password_reset", data)
}

// SendWelcomeEmail sends a welcome email to new users
func (s *Service) SendWelcomeEmail(to, displayName string) error {
	data := map[string]interface{}{
		"DisplayName": displayName,
	}

	subject := "Welcome to NexusFlow!"
	return s.sendEmail(to, subject, "welcome", data)
}

// sendEmail sends an email using a template
func (s *Service) sendEmail(to, subject, templateName string, data map[string]interface{}) error {
	// Get template
	tmpl, ok := s.templates[templateName]
	if !ok {
		return fmt.Errorf("template not found: %s", templateName)
	}

	// Render template
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Create message
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromAddress))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	// Send email
	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// Inline HTML templates
const verificationTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify Your Email</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #3b82f6 0%, #6366f1 100%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
        .content { background: #f9fafb; padding: 30px; border-radius: 0 0 10px 10px; }
        .button { display: inline-block; padding: 12px 30px; background: #3b82f6; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 30px; color: #6b7280; font-size: 14px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🚀 NexusFlow</h1>
    </div>
    <div class="content">
        <h2>Hi {{.DisplayName}}!</h2>
        <p>Thanks for signing up for NexusFlow. Please verify your email address to get started.</p>
        <p style="text-align: center;">
            <a href="{{.VerificationURL}}" class="button">Verify Email Address</a>
        </p>
        <p>Or copy and paste this link into your browser:</p>
        <p style="word-break: break-all; color: #3b82f6;">{{.VerificationURL}}</p>
        <p>This link will expire in 24 hours.</p>
    </div>
    <div class="footer">
        <p>If you didn't create an account, you can safely ignore this email.</p>
        <p>&copy; 2024 NexusFlow. All rights reserved.</p>
    </div>
</body>
</html>
`

const passwordResetTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Your Password</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #3b82f6 0%, #6366f1 100%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
        .content { background: #f9fafb; padding: 30px; border-radius: 0 0 10px 10px; }
        .button { display: inline-block; padding: 12px 30px; background: #3b82f6; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 30px; color: #6b7280; font-size: 14px; }
        .warning { background: #fef3c7; border-left: 4px solid #f59e0b; padding: 15px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🔒 NexusFlow</h1>
    </div>
    <div class="content">
        <h2>Hi {{.DisplayName}}!</h2>
        <p>We received a request to reset your password. Click the button below to create a new password:</p>
        <p style="text-align: center;">
            <a href="{{.ResetURL}}" class="button">Reset Password</a>
        </p>
        <p>Or copy and paste this link into your browser:</p>
        <p style="word-break: break-all; color: #3b82f6;">{{.ResetURL}}</p>
        <div class="warning">
            <strong>⚠️ Security Notice:</strong> This link will expire in 24 hours. If you didn't request a password reset, please ignore this email and your password will remain unchanged.
        </div>
    </div>
    <div class="footer">
        <p>For security reasons, never share this link with anyone.</p>
        <p>&copy; 2024 NexusFlow. All rights reserved.</p>
    </div>
</body>
</html>
`

const welcomeTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to NexusFlow</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #3b82f6 0%, #6366f1 100%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
        .content { background: #f9fafb; padding: 30px; border-radius: 0 0 10px 10px; }
        .feature { background: white; padding: 15px; margin: 15px 0; border-radius: 5px; border-left: 4px solid #3b82f6; }
        .footer { text-align: center; margin-top: 30px; color: #6b7280; font-size: 14px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🎉 Welcome to NexusFlow!</h1>
    </div>
    <div class="content">
        <h2>Hi {{.DisplayName}}!</h2>
        <p>Welcome aboard! We're excited to have you join NexusFlow, the open-source project management platform you can own.</p>
        
        <h3>What's Next?</h3>
        <div class="feature">
            <strong>📋 Create Your First Project</strong><br>
            Start organizing your work with projects, boards, and issues.
        </div>
        <div class="feature">
            <strong>👥 Invite Your Team</strong><br>
            Collaborate with your colleagues and manage permissions.
        </div>
        <div class="feature">
            <strong>⚙️ Customize Your Workspace</strong><br>
            Tailor workflows, fields, and settings to match your needs.
        </div>
        
        <p>Need help getting started? Check out our <a href="https://docs.nexusflow.io" style="color: #3b82f6;">documentation</a> or reach out to our support team.</p>
    </div>
    <div class="footer">
        <p>Happy project managing!</p>
        <p>&copy; 2024 NexusFlow. All rights reserved.</p>
    </div>
</body>
</html>
`
