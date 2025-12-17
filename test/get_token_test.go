package test

import (
	"cashly/internal/entity"
	tokenservice "cashly/internal/service/token"
	"cashly/internal/service/token/mocks"
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetBankToken(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	const (
		familyID = 1
		userID   = int64(42)
	)

	tests := []struct {
		name             string
		mockToken        *entity.UserBankToken
		mockRepoErr      error
		mockDecryptedVal string
		mockDecryptErr   error
		wantHasToken     bool
		wantToken        *entity.UserBankToken
		wantErr          bool
	}{
		{
			name: "user has token - successfully decrypted",
			mockToken: &entity.UserBankToken{
				ID:        1,
				UserID:    userID,
				FamilyID:  familyID,
				Token:     "encrypted_token",
				CreatedAt: fixedTime,
			},
			mockDecryptedVal: "decrypted_token",
			wantHasToken:     true,
			wantToken: &entity.UserBankToken{
				ID:        1,
				UserID:    userID,
				FamilyID:  familyID,
				Token:     "decrypted_token",
				CreatedAt: fixedTime,
			},
		},
		{
			name:         "user has no token",
			mockRepoErr:  pgx.ErrNoRows,
			wantHasToken: false,
			wantToken:    nil,
		},
		{
			name: "decryption fails",
			mockToken: &entity.UserBankToken{
				ID:    1,
				Token: "encrypted_token",
			},
			mockDecryptErr: assert.AnError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.TokenIfaceMock)
			mockEncryptor := new(mocks.EncryptorMock)

			mockService.On("Get", context.Background(), familyID, userID).
				Return(tt.mockToken, tt.mockRepoErr)

			if tt.mockToken != nil {
				mockEncryptor.On("Decrypt", tt.mockToken.Token).
					Return(tt.mockDecryptedVal, tt.mockDecryptErr)
			}

			svc := tokenservice.New(mockEncryptor, mockService, newTestLogger())

			hasToken, got, err := svc.Get(context.Background(), familyID, userID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantHasToken, hasToken)
			assert.Equal(t, tt.wantToken, got)

			mockService.AssertExpectations(t)
			mockEncryptor.AssertExpectations(t)
		})
	}
}
