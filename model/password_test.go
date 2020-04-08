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
	salt, err := createSalt()
	assert.NoError(t, err)
	assert.NotNil(t, salt)

}

func TestEncodePassword(t *testing.T) {
	correctPassword := "$argon2id$v=19$m=65536,t=3,p=1$L/cXOPSeeE9f68JKienFug$t7/BEg7jx2Gjx/OwSFyWgw"
	c := NewPasswordConfig()
	passByte := []byte("ostap")
	b64Salt := "L/cXOPSeeE9f68JKienFug"

	salt, err := base64.RawStdEncoding.DecodeString(b64Salt)
	assert.NoError(t, err)
	hash := argon2.IDKey(passByte, salt, c.time, c.memory, c.threads, c.keyLen)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	good := fmt.Sprintf(hashFormat, argon2.Version, c.memory, c.time, c.threads, b64Salt, b64Hash)
	assert.Equal(t, correctPassword, good)

	saltNew, err := createSalt()
	assert.NoError(t, err)

	bSalt := base64.RawStdEncoding.EncodeToString(saltNew)
	hash1 := argon2.IDKey(passByte, saltNew, c.time, c.memory, c.threads, c.keyLen)
	b64Hash1 := base64.RawStdEncoding.EncodeToString(hash1)
	wrong := fmt.Sprintf(hashFormat, argon2.Version, c.memory, c.time, c.threads, bSalt, b64Hash1)
	assert.NotEqual(t, correctPassword, wrong)
	fmt.Println(wrong)
}

func TestComparePassword(t *testing.T) {
	correctPassword := "$argon2id$v=19$m=65536,t=3,p=1$gv1q09I+VqtsT64dGOClcQ$tM+aG4UJ3d5xAf5smeY/3A"
	passString := "ostap"
	passWrong := "wrongPassword"

	isSame, err := ComparePassword(passString, correctPassword)
	assert.True(t, isSame)
	assert.NoError(t, err)

	isBad, err := ComparePassword(passWrong, correctPassword)
	assert.False(t, isBad)
	assert.NoError(t, err)
}
