package test

import (
	"cashly/internal/entity"
	userservice "cashly/internal/service/user"
	"cashly/internal/service/user/mocks"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	const uid = int64(42)

	tests := []struct {
		name        string
		mockRepoErr error
		wantUser    *entity.User
		wantErr     bool
	}{
		{
			name: "success",
			wantUser: &entity.User{
				ID:        42,
				Username:  "username1",
				Firstname: "firstname1",
				JoinedAt:  fixedTime,
			},
			mockRepoErr: nil,
			wantErr:     false,
		},
		{
			name:        "database error",
			wantUser:    nil,
			mockRepoErr: assert.AnError,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.UserIfaceMock)

			mockService.On("GetByID", context.Background(), uid).Return(tt.wantUser, tt.mockRepoErr)

			svc := userservice.New(mockService, nil, "", nil, newTestLogger())

			gotUser, err := svc.GetByID(context.Background(), uid)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.wantUser, gotUser)
		})
	}
}
