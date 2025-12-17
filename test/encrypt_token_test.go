package test

import (
	tokenservice "cashly/internal/service/token"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncryption(t *testing.T) {
	key := []byte("examplekey1234567890examplekey12")

	te := tokenservice.NewEncrypt(key)

	original := "sensitive_bank_token"

	encrypted, err := te.Encrypt(original)
	require.NoError(t, err)

	decrypted, err := te.Decrypt(encrypted)
	require.NoError(t, err)

	require.Equal(t, original, decrypted)
}
