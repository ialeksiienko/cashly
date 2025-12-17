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

func TestRegisterUser(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	const (
		familyID = 1
		userID   = int64(42)
	)

	tests := []struct {
		name        string
		mockUser    *entity.User
		mockRepoErr error
		wantUser    *entity.User
		wantErr     bool
	}{
		{
			name: "success",
			mockUser: &entity.User{
				ID:        1,
				Username:  "username",
				Firstname: "firstname",
			},
			mockRepoErr: nil,
			wantUser: &entity.User{
				ID:        1,
				Username:  "username",
				Firstname: "firstname",
				JoinedAt:  fixedTime,
			},
			wantErr: false,
		},
		{
			name: "database error",
			mockUser: &entity.User{
				ID:        1,
				Username:  "username",
				Firstname: "firstname",
			},
			mockRepoErr: assert.AnError,
			wantUser:    nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.UserIfaceMock)

			mockService.On("Save", context.Background(), tt.mockUser).Return(tt.wantUser, tt.mockRepoErr)

			svc := userservice.New(mockService, nil, "", nil, newTestLogger())

			gotUser, err := svc.Register(context.Background(), tt.mockUser)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.wantUser, gotUser)
		})
	}
}
