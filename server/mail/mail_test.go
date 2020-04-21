package mail

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestGenerateVerificationCode(t *testing.T) {
	code := GenerateVerificationCode()

	assert.Equal(t, activationCodeSize, len(code))
}

func TestCreateEmail(t *testing.T) {
	result := "Hello! Hello! Hello!"
	email := &EmailInfo{
		Sender:    "user@example.com",
		Password:  "123456789",
		Host:      "localhost",
		Port:      8099,
		Recipient: "recipient@example.com",
		Subject:   "Testing 1",
		Body:      result,
		URL:       "i.try.send.email.com",
		Code:      "1234567890qwertyuiop",
	}

	mail := email.CreateEmail(result)

	assert.Equal(t, email.Sender, mail.GetHeader("From")[0])
	assert.Equal(t, email.Recipient, mail.GetHeader("To")[0])
	assert.Equal(t, email.Subject, mail.GetHeader("Subject")[0])
	assert.Equal(t, email.Body, result)
}

func TestSetupURLQueryParameters(t *testing.T) {
	result := "Hello! Hello! Hello!"
	email := &EmailInfo{
		Sender:    "user@example.com",
		Password:  "123456789",
		Host:      "localhost",
		Port:      8099,
		Recipient: "recipient@example.com",
		Subject:   "Testing 1",
		Body:      result,
		URL:       "i.try.send.email.com",
		Code:      "1234567890qwertyuiop",
	}

	urlQuery, err := email.SetupURLQueryParameters("123456789", "bodja")

	require.NoError(t, err)
	assert.Equal(t, "i.try.send.email.com?code=123456789&login=bodja", urlQuery)
	assert.NotEqual(t, "i.try.send.email.com?code=qwertyui&login=12345", urlQuery)
}
