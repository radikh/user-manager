package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/lvl484/user-manager/config"

	"github.com/lvl484/user-manager/server/mail"
)

const (
	queryAddActivationCode          = `INSERT INTO email_codes(id, user_name, email, verification_code, created_at) VALUES ($1,$2,$3,$4,$5)`
	queryDeleteActivationCode       = `DELETE FROM email_codes WHERE user_name=$1`
	querySelectVerificationCodeTime = `SELECT verification_code, created_at FROM email_codes WHERE user_name=$1`
)

type Verification struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

// AddActivationCode adds new activation code for user to database
func (ur *UsersRepo) AddActivationCode(user *User) error {
	emailConfig, err := SetupEmailComponents(user.Email)
	if err != nil {
		return err
	}

	verificationCode := mail.GenerateVerificationCode()

	err = emailConfig.SendMail(user.Username, verificationCode)
	if err != nil {
		return fmt.Errorf("AddActivationCode SendMail error: %w", err)
	}

	_, err = ur.db.Exec(queryAddActivationCode, user.ID, user.Username, user.Email, verificationCode, time.Now())
	return err
}

// DeleteVerificationCode deletes activation code for user from database
func (ur *UsersRepo) DeleteVerificationCode(login string) error {
	_, err := ur.db.Exec(queryDeleteActivationCode, login)

	return err
}

// GetVerificationCodeTime gets verification code and time when it was code was created from database
func (ur *UsersRepo) GetVerificationCodeTime(login string) (*time.Time, string, error) {
	var verifyCodeTime *time.Time
	var verificationCode string

	err := ur.db.QueryRow(querySelectVerificationCodeTime, login).Scan(&verificationCode, &verifyCodeTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", errors.Wrap(err, msgUserDidNotExist)
		}
		return nil, "", err
	}

	return verifyCodeTime, verificationCode, nil
}

// SetupEmailComponents setups email components
func SetupEmailComponents(email string) (*mail.EmailInfo, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("SetupEmailComponents NewConfig error: %w", err)
	}

	emailConfig, err := cfg.EmailConfig()
	if err != nil {
		return nil, fmt.Errorf("SetupEmailComponents EmailConfig error: %w", err)
	}

	return &mail.EmailInfo{
		Sender:    emailConfig.Sender,
		Password:  emailConfig.Password,
		Host:      emailConfig.Host,
		Port:      emailConfig.Port,
		Recipient: email,
		Subject:   mail.EmailSubject,
		Body:      mail.EmailBody,
		URL:       mail.EmailContentLink,
	}, nil
}
