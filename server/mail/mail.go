package mail

import (
	"bytes"
	"html/template"
	"log"
	"math/rand"
	"net/url"
	"time"

	"github.com/lvl484/user-manager/logger"

	"gopkg.in/gomail.v2"
)

const (
	digitsForActivationCode          = "0123456789"
	specialsSymbolsForActivationCode = "~=+%^*/()[]{}/!@#$?|"
	literalsSymbolsForActivationCode = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	activationCodeSize = 50

	htmlTemplateName = "template.html"
	htmlTemplatePath = "server/mail/mail_template/template.html"

	EmailSubject = "Testing e-mail!"
	EmailBody    = "This is email body. \r\n If you did not receive email letter, please looking for in the SPAM."
	//EmailContentLink = "http://127.0.0.1:8000/verification"
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

// GenerateVerificationCode created random string with length = activationCodeSize
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
func (mc EmailInfo) SendMail(login string, code string) {
	mc.Code = code

	mc.URL = setupURLQueryParameters(mc, code, login)

	result := formattingByTemplate(&mc)

	email := createEmail(&mc, result)

	dialer := gomail.NewDialer(mc.Host, mc.Port, mc.Sender, mc.Password)

	// Send the email to Recipient
	if err := dialer.DialAndSend(email); err != nil {
		logger.LogUM.Fatalf("Send email error: %v", err)
	}
}

func setupURLQueryParameters(mc EmailInfo, code string, login string) string {
	u, err := url.Parse(mc.URL)
	if err != nil {
		logger.LogUM.Fatalf("parse url error: %v", err)
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		logger.LogUM.Fatalf("parse query error: %v", err)
	}

	q.Set("code", code)
	q.Set("login", login)

	u.RawQuery = q.Encode()

	return u.String()
}

func formattingByTemplate(emailInfo *EmailInfo) string {
	t := template.New(htmlTemplateName)

	t, err := t.ParseFiles(htmlTemplatePath)
	if err != nil {
		log.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, emailInfo); err != nil {
		log.Println(err)
	}

	return tpl.String()
}

// createEmail create email letter structure
func createEmail(emailInfo *EmailInfo, result string) *gomail.Message {
	email := gomail.NewMessage()

	email.SetHeader("From", emailInfo.Sender)
	email.SetHeader("To", emailInfo.Recipient)
	//email.SetAddressHeader("Cc", "<RECIPIENT CC>", "<RECIPIENT CC NAME>")
	email.SetHeader("Subject", emailInfo.Subject)
	email.SetBody("text/html", result)
	//email.Attach("/home/bodja/Desktop/UM/user-manager/server/mail/

	return email
}
