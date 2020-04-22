package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/url"
	"path"

	"gopkg.in/gomail.v2"
)

const (
	verificationTemplateName = "verification_email.html"

	verificationEmailSubject = "Registration UM service!"
)

// EmailInfo includes necessary fields, which need for email authorization
type EmailConfig struct {
	Sender       string
	Password     string
	Host         string
	Port         int
	PublicURL    string
	TemplatePath string
}

type VerificationEmail struct {
	UserEmail string
	Code      string
	Username  string
}

// VerificationTemplate includes fields which are using for verification page template
type VerificationTemplate struct {
	VerificationEmail
	URL string
}

func (mc *EmailConfig) SendVerificationMail(verification *VerificationEmail) error {
	verURL, err := generateVerificationURL(mc.PublicURL, verification.Code, verification.Username)
	if err != nil {
		return fmt.Errorf("SendVerificationMail generateVerificationURL error: %w", err)
	}

	verificationTemplate := &VerificationTemplate{
		VerificationEmail: *verification,
		URL:               verURL.String(),
	}

	templatePath := path.Join(mc.TemplatePath, verificationTemplateName)

	result, err := renderTemplate(templatePath, verificationTemplate)
	if err != nil {
		return fmt.Errorf("SendVerificationMail renderTemplate error: %w", err)
	}

	email := gomail.NewMessage()

	email.SetHeader("From", mc.Sender)
	email.SetHeader("To", verification.UserEmail)
	email.SetHeader("Subject", verificationEmailSubject)
	email.SetBody("text/html", result)

	dialer := gomail.NewDialer(mc.Host, mc.Port, mc.Sender, mc.Password)

	// Send the email to Recipient
	if err := dialer.DialAndSend(email); err != nil {
		return fmt.Errorf("SendMail DialAndSend error: %w", err)
	}

	return nil
}

func generateVerificationURL(publicURL, code, login string) (*url.URL, error) {
	u, err := url.Parse(publicURL + "/verification")
	if err != nil {
		return nil, fmt.Errorf("generateVerificationURL parse url error: %w", err)
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("generateVerificationURL parse url query error: %w", err)
	}

	q.Set("code", code)
	q.Set("login", login)

	u.RawQuery = q.Encode()

	return u, nil
}

func renderTemplate(templateName string, params interface{}) (string, error) {
	t, err := template.ParseFiles(templateName)
	if err != nil {
		return "", fmt.Errorf("renderTemplate ParseFiles error: %w", err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, params); err != nil {
		log.Println(err)
		return "", fmt.Errorf("renderTemplate Execute error: %w", err)
	}

	return tpl.String(), nil
}
