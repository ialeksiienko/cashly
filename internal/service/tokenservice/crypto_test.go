package tokenservice

import (
	"cashly/internal/pkg/sl"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt(t *testing.T) {
	var key [32]byte
	copy(key[:], []byte("examplekey1234567890examplekey12"))

	te := NewEncrypt(key, &sl.MyLogger{})

	original := "sensitive_bank_token"

	encrypted, err := te.Encrypt(original)
	require.NoError(t, err)

	decrypted, err := te.Decrypt(encrypted)
	require.NoError(t, err)
	require.Equal(t, original, decrypted)
}
