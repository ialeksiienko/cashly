package test

import (
	"cashly/internal/entity"
	tokenservice "cashly/internal/service/token"
	"cashly/internal/service/token/mocks"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSaveBankToken(t *testing.T) {
	mockService := new(mocks.TokenIfaceMock)
	mockEncryptor := new(mocks.EncryptorMock)

	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	expectedUBT := &entity.UserBankToken{
		ID:        1,
		UserID:    42,
		FamilyID:  1,
		Token:     "secret_token",
		CreatedAt: fixedTime,
	}

	mockService.On("Save", context.Background(), 1, int64(42), "secret_token").Return(expectedUBT, nil)

	svc := tokenservice.New(mockEncryptor, mockService, newTestLogger())

	got, err := svc.Save(context.Background(), 1, int64(42), "secret_token")

	assert.NoError(t, err)
	assert.Equal(t, expectedUBT, got)

	mockService.AssertExpectations(t)
}
