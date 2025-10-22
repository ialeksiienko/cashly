package tokenservice

import (
	"cashly/internal/pkg/sl"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
)

type Encrypt struct {
	key [32]byte
	sl  sl.Logger
}

func NewEncrypt(key [32]byte, sl sl.Logger) *Encrypt {
	return &Encrypt{key: key, sl: sl}
}

func (te *Encrypt) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(te.key[:])
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)

	final := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(final), nil
}

func (te *Encrypt) Decrypt(encrypted string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(te.key[:])
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("incorrect encrypted string")
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
