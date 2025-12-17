package test

import (
	"cashly/internal/entity"
	familyservice "cashly/internal/service/family"
	"cashly/internal/service/family/mocks"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetFamilyByCode(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name          string
		mockRepoErr   error
		mockFamily    *entity.Family
		mockExpiresAt time.Time
		wantErr       bool
	}{
		{
			name:        "success",
			mockRepoErr: nil,
			mockFamily: &entity.Family{
				ID:        1,
				CreatedBy: 42,
				Name:      "Test Family",
			},
			mockExpiresAt: fixedTime,
			wantErr:       false,
		},
		{
			name:          "database error",
			mockRepoErr:   assert.AnError,
			mockFamily:    nil,
			mockExpiresAt: time.Time{},
			wantErr:       true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.FamilyIfaceMock)

			mockService.On("GetByCode", context.Background(), "ABC123").
				Return(tt.mockFamily, tt.mockExpiresAt, tt.mockRepoErr)

			svc := familyservice.New(mockService, nil, newTestLogger())

			gotFamily, gotExpiresAt, err := svc.GetByCode(context.Background(), "ABC123")

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.mockFamily, gotFamily)
			assert.Equal(t, tt.mockExpiresAt, gotExpiresAt)

			mockService.AssertExpectations(t)
		})
	}
}
