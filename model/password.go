// Package model provides user-manager specific data structures,
// which are meant to be used across the whole application.
package model

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

const (
	configTime            = 3
	configMemory          = 64 * 1024
	configThreads         = 1
	configKeyLen          = 16
	lengthSalt            = 16
	hashFormat            = "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	hashSplit             = "$"
	configFormat          = "m=%d,t=%d,p=%d"
	messageSaltError      = "Error generating salt for password"
	messageConfigDecode   = "Decoding of password's configuration failed"
	messageSaltDecode     = "Decoding of password's salt failed"
	messagePasswordDecode = "Decoding of user's password failed"
)

// PasswordConfig is structure that describes complication of hashing the password
type PasswordConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

// NewPasswordConfig returns config for encode
func NewPasswordConfig() *PasswordConfig {
	return &PasswordConfig{
		time:    configTime,
		memory:  configMemory,
		threads: configThreads,
		keyLen:  configKeyLen,
	}
}

// createSalt create random salt according to lengthSalt
func createSalt() ([]byte, error) {
	salt := make([]byte, lengthSalt)
	_, err := rand.Read(salt)
	return salt, err
}

// EncodePassword returns encoded password
func EncodePassword(c *PasswordConfig, pass string) (string, error) {
	salt, err := createSalt()
	if err != nil {
		return "", errors.Wrap(err, messageSaltError)
	}

	passByte := []byte(pass)

	hash := argon2.IDKey(passByte, salt, c.time, c.memory, c.threads, c.keyLen)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	full := fmt.Sprintf(hashFormat, argon2.Version, c.memory, c.time, c.threads, b64Salt, b64Hash)

	return full, nil
}

// ComparePassword returns true if password matches
func ComparePassword(password, hash string) (bool, error) {
	parts := strings.Split(hash, hashSplit)
	c := &PasswordConfig{}
	_, err := fmt.Sscanf(parts[3], configFormat, &c.memory, &c.time, &c.threads)

	if err != nil {
		return false, errors.Wrap(err, messageConfigDecode)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])

	if err != nil {
		return false, errors.Wrap(err, messageSaltDecode)
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])

	if err != nil {
		return false, errors.Wrap(err, messagePasswordDecode)
	}

	c.keyLen = uint32(len(decodedHash))
	comparisonHash := argon2.IDKey([]byte(password), salt, c.time, c.memory, c.threads, c.keyLen)

	return (subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1), nil
}
