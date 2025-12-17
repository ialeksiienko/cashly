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

func TestGetUsersByFamilyID(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	const fid = 1

	tests := []struct {
		name        string
		mockRepoErr error
		wantUsers   []entity.User
		wantErr     bool
	}{
		{
			name: "success",
			wantUsers: []entity.User{
				{
					ID:        42,
					Username:  "username1",
					Firstname: "firstname1",
					JoinedAt:  fixedTime,
				},
			},
			mockRepoErr: nil,
			wantErr:     false,
		},
		{
			name:        "database error",
			wantUsers:   nil,
			mockRepoErr: assert.AnError,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.UserIfaceMock)

			mockService.On("GetAllUsersInFamily", context.Background(), fid).Return(tt.wantUsers, tt.mockRepoErr)

			svc := userservice.New(mockService, nil, "", nil, newTestLogger())

			gotUser, err := svc.GetUsersByFamilyID(context.Background(), fid)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.wantUsers, gotUser)
		})
	}
}
