package familyservice

import (
	"crypto/rand"
	"encoding/base32"
	"strings"
)

const (
	codeLength = 6
)

func (s Service) GenerateInviteCode() (string, error) {
	b := make([]byte, 5)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	code := base32.StdEncoding.EncodeToString(b)
	code = strings.ToUpper(code)
	return code[:codeLength], nil
}
