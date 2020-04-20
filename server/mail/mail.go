package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/url"
	"time"

	"gopkg.in/gomail.v2"
)

const (
	digitsForActivationCode          = "0123456789"
	specialsSymbolsForActivationCode = "~=+%^*/()[]{}/!@#$?|"
	literalsSymbolsForActivationCode = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	activationCodeSize = 50

	htmlTemplateName = "template.html"
	htmlTemplatePath = "server/mail/mail_template/template.html"

	EmailSubject     = "Testing e-mail!"
	EmailBody        = "This is email body. \r\n If you did not receive email letter, please looking for in the SPAM."
	EmailContentLink = "http://127.0.0.1:8000/verification"
)

// EmailInfo includes necessary fields, which need for email authorization
type EmailInfo struct {
	Sender    string
	Password  string
	Host      string
	Port      int
	Recipient string
	Subject   string
	Body      string
	URL       string
	Code      string
}

// GenerateVerificationCode creates random string with length = activationCodeSize
func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())

	all := literalsSymbolsForActivationCode + digitsForActivationCode + specialsSymbolsForActivationCode

	buf := make([]byte, activationCodeSize)

	buf[0] = digitsForActivationCode[rand.Intn(len(digitsForActivationCode))]
	buf[1] = specialsSymbolsForActivationCode[rand.Intn(len(specialsSymbolsForActivationCode))]

	for i := 2; i < activationCodeSize; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}

	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})

	return string(buf)
}

// SendMail uses prepared template, parses it and converts into string.
// Set up necessary fields like Sender, Recipient, Subject, email Body.
// Send created email to recipient.
func (mc EmailInfo) SendMail(login string, code string) error {
	mc.Code = code

	var err error
	mc.URL, err = setupURLQueryParameters(mc, code, login)
	if err != nil {
		return fmt.Errorf("SendMail setupURLQueryParameters error: %w", err)
	}

	result, err := formattingByTemplate(&mc)
	if err != nil {
		return fmt.Errorf("SendMail formattingByTemplate error: %w", err)
	}

	email := createEmail(&mc, result)

	dialer := gomail.NewDialer(mc.Host, mc.Port, mc.Sender, mc.Password)

	// Send the email to Recipient
	if err := dialer.DialAndSend(email); err != nil {
		return fmt.Errorf("SendMail DialAndSend error: %w", err)
	}

	return nil
}

func setupURLQueryParameters(mc EmailInfo, code string, login string) (string, error) {
	u, err := url.Parse(mc.URL)
	if err != nil {
		return "", fmt.Errorf("setupURLQueryParameters parse url error: %w", err)
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", fmt.Errorf("setupURLQueryParameters parse url query error: %w", err)
	}

	q.Set("code", code)
	q.Set("login", login)

	u.RawQuery = q.Encode()

	return u.String(), nil
}

func formattingByTemplate(emailInfo *EmailInfo) (string, error) {
	t := template.New(htmlTemplateName)

	t, err := t.ParseFiles(htmlTemplatePath)
	if err != nil {
		return "", fmt.Errorf("formattingByTemplate ParseFiles error: %w", err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, emailInfo); err != nil {
		log.Println(err)
		return "", fmt.Errorf("formattingByTemplate Execute error: %w", err)
	}

	return tpl.String(), nil
}

// createEmail create email letter structure
func createEmail(emailInfo *EmailInfo, result string) *gomail.Message {
	email := gomail.NewMessage()

	email.SetHeader("From", emailInfo.Sender)
	email.SetHeader("To", emailInfo.Recipient)
	email.SetHeader("Subject", emailInfo.Subject)
	email.SetBody("text/html", result)

	return email
}
