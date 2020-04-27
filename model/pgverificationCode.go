package model

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/pkg/errors"

	"github.com/lvl484/user-manager/config"

	"github.com/lvl484/user-manager/mail"
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
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("AddActivationCode NewConfig error: %w", err)
	}

	emailConfig, err := cfg.EmailConfig()
	if err != nil {
		return fmt.Errorf("AddActivationCode EmailConfig error: %w", err)
	}

	verificationEmail := &mail.VerificationEmail{
		UserEmail: user.Email,
		Code:      generateVerificationCode(),
		Username:  user.Username,
	}

	err = emailConfig.SendVerificationMail(verificationEmail)
	if err != nil {
		return fmt.Errorf("AddActivationCode SendVerificationMail error: %w", err)
	}

	_, err = ur.db.Exec(queryAddActivationCode, user.ID, verificationEmail.Username, verificationEmail.UserEmail, verificationEmail.Code, time.Now())
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

// generateVerificationCode creates random string with length = activationCodeSize
func generateVerificationCode() string {
	const (
		digitsForActivationCode          = "0123456789"
		specialsSymbolsForActivationCode = "~=+%^*/()[]{}/!@#$?|"
		literalsSymbolsForActivationCode = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
		activationCodeSize               = 50
	)

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
