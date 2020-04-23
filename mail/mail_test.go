package mail

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestSendMail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pass := os.Getenv("EMAIL_PASSWORD")
	if pass == "" {
		t.Skip("EMAIL_PASSWORD is not set... skipping integration test")
	}

	mail := &EmailConfig{
		Sender:       "user.namager@gmail.com",
		Password:     pass,
		Host:         "smtp.gmail.com",
		Port:         587,
		TemplatePath: "mail_template/",
		PublicURL:    "http://localhost:8000",
	}
	vm := &VerificationEmail{
		UserEmail: "user1@example.com",
		Code:      "123456789",
		Username:  "User1",
	}
	err := mail.SendVerificationMail(vm)
	assert.NoError(t, err)
}

func TestGenerateVerificationURLError(t *testing.T) {
	urlQuery, err := generateVerificationURL("^%", "123456789", "bodja")
	require.Error(t, err)

	assert.Nil(t, urlQuery)
	assert.Error(t, err)
}

func TestGenerateVerificationURL(t *testing.T) {
	urlQuery, err := generateVerificationURL("http://localhost:8000", "123456789", "bodja")
	require.NoError(t, err)

	assert.Equal(t, "http://localhost:8000/verification?code=123456789&login=bodja", urlQuery.String())
}

func TestRenderTemplateWrongTemplatePath(t *testing.T) {
	emailTemplate, err := renderTemplate("wrong/path", nil)
	require.Error(t, err)

	assert.Empty(t, emailTemplate)
}
