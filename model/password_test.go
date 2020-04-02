// Package model provides user-manager specific data structures,
// which are meant to be used across the whole application.
package model

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/argon2"
)

func TestNewPasswordConfig(t *testing.T) {

	goodConfig := &PasswordConfig{
		time:    configTime,
		memory:  configMemory,
		threads: configThreads,
		keyLen:  configKeyLen,
	}
	badConfig := &PasswordConfig{
		time:    5,
		memory:  configMemory,
		threads: configThreads,
		keyLen:  4,
	}
	got := NewPasswordConfig()
	assert.Equal(t, goodConfig, got)
	assert.NotEqual(t, badConfig, got)

}

func Test_createSalt(t *testing.T) {
	_, err := createSalt()
	assert.NoError(t, err)
}

func TestEncodePassword(t *testing.T) {
	correctPassword := "$argon2id$v=19$m=65536,t=3,p=1$BCDndJ1kUOAAW/mwP7ViOQ$Ig4hpteBW1YM7Lrh3EHkHQ"
	c := NewPasswordConfig()
	passByte := []byte("password")
	b64Salt := "BCDndJ1kUOAAW/mwP7ViOQ"

	salt, _ := base64.RawStdEncoding.DecodeString(b64Salt)
	hash := argon2.IDKey(passByte, salt, c.time, c.memory, c.threads, c.keyLen)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	saltNew, _ := createSalt()
	bSalt := base64.RawStdEncoding.EncodeToString(saltNew)

	good := fmt.Sprintf(hashFormat, argon2.Version, c.memory, c.time, c.threads, b64Salt, b64Hash)
	assert.Equal(t, correctPassword, good)

	wrong := fmt.Sprintf(hashFormat, argon2.Version, c.memory, c.time, c.threads, bSalt, b64Hash)
	assert.NotEqual(t, correctPassword, wrong)
}

func TestComparePassword(t *testing.T) {
	correctPassword := "$argon2id$v=19$m=65536,t=3,p=1$BCDndJ1kUOAAW/mwP7ViOQ$Ig4hpteBW1YM7Lrh3EHkHQ"
	passString := "password"
	passWrong := "wrongPassword"

	isSame, err := ComparePassword(passString, correctPassword)
	assert.True(t, isSame)
	assert.NoError(t, err)

	isBad, err := ComparePassword(passWrong, correctPassword)
	assert.False(t, isBad)
	assert.NoError(t, err)
}
