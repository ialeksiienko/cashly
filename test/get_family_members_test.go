package test

import (
	"cashly/internal/entity"
	userservice "cashly/internal/service/user"
	"cashly/internal/service/user/mocks"
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetFamilyMembers(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	const (
		familyID = 1
		userID   = int64(42)
	)

	tests := []struct {
		name         string
		mockRepoErr  error
		mockTokenErr error
		wantUsers    []entity.User
		wantMembers  []userservice.Member
		wantUBT      *entity.UserBankToken
		wantHasToken bool
		wantErr      bool
	}{
		{
			name:         "success",
			mockRepoErr:  nil,
			mockTokenErr: nil,
			wantUsers: []entity.User{
				{
					ID:        42,
					Username:  "username1",
					Firstname: "firstname1",
					JoinedAt:  fixedTime,
				},
			},
			wantMembers: []userservice.Member{
				{
					ID:        42,
					Username:  "username1",
					Firstname: "firstname1",
					IsAdmin:   true,
					IsCurrent: true,
					HasToken:  true,
				},
			},
			wantUBT: &entity.UserBankToken{
				ID:        1,
				UserID:    userID,
				FamilyID:  familyID,
				Token:     "encrypted_token",
				CreatedAt: fixedTime,
			},
			wantHasToken: true,
			wantErr:      false,
		},
		{
			name:         "no members in family",
			mockRepoErr:  assert.AnError,
			mockTokenErr: nil,
			wantUsers:    nil,
			wantMembers:  nil,
			wantUBT:      nil,
			wantHasToken: false,
			wantErr:      true,
		},
		{
			name:         "member has no token",
			mockRepoErr:  nil,
			mockTokenErr: pgx.ErrNoRows,
			wantUsers: []entity.User{
				{
					ID:        42,
					Username:  "username1",
					Firstname: "firstname1",
					JoinedAt:  fixedTime,
				},
			},
			wantMembers: []userservice.Member{
				{
					ID:        42,
					Username:  "username1",
					Firstname: "firstname1",
					IsAdmin:   true,
					IsCurrent: true,
					HasToken:  false,
				},
			},
			wantUBT:      nil,
			wantHasToken: false,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.UserIfaceMock)
			mockToken := new(mocks.TokenProviderMock)

			mockService.On("GetAllUsersInFamily", context.Background(), familyID).Return(tt.wantUsers, tt.mockRepoErr)
			mockToken.On("Get", context.Background(), familyID, userID).
				Return(tt.wantHasToken, tt.wantUBT, tt.mockTokenErr)

			srv := userservice.New(mockService, nil, "", mockToken, newTestLogger())

			gotMembers, err := srv.GetFamilyMembers(context.Background(), &entity.Family{
				ID:        familyID,
				CreatedBy: userID,
				Name:      "test",
			}, userID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.wantMembers, gotMembers)
			t.Logf("%+v", gotMembers)
		})
	}
}
