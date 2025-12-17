package test

import (
	userservice "cashly/internal/service/user"
	"cashly/internal/service/user/mocks"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveUserToFamily(t *testing.T) {
	const (
		familyID = 1
		userID   = int64(42)
	)

	tests := []struct {
		name        string
		mockRepoErr error
		wantErr     bool
	}{
		{
			name:        "success",
			mockRepoErr: nil,
			wantErr:     false,
		},
		{
			name:        "database error",
			mockRepoErr: assert.AnError,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.UserIfaceMock)

			mockService.On("SaveToFamily", context.Background(), familyID, userID).Return(tt.mockRepoErr)

			svc := userservice.New(mockService, nil, "", nil, newTestLogger())

			gotErr := svc.SaveToFamily(context.Background(), familyID, userID)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}
		})
	}
}
