package test

import (
	"cashly/internal/entity"
	familyservice "cashly/internal/service/family"
	"cashly/internal/service/family/mocks"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFamiliesByUserID(t *testing.T) {
	tests := []struct {
		name         string
		mockRepoErr  error
		mockFamilies []entity.Family
		wantErr      bool
	}{
		{
			name:        "success",
			mockRepoErr: nil,
			mockFamilies: []entity.Family{
				{
					ID:        1,
					CreatedBy: 42,
					Name:      "Test Family",
				},
			},
			wantErr: false,
		},
		{
			name:         "database error",
			mockRepoErr:  nil,
			mockFamilies: nil,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.FamilyIfaceMock)

			mockService.On("GetFamiliesByUserID", context.Background(), int64(42)).
				Return(tt.mockFamilies, tt.mockRepoErr)

			svc := familyservice.New(mockService, nil, newTestLogger())

			got, err := svc.GetFamiliesByUserID(context.Background(), int64(42))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.mockFamilies, got)

			mockService.AssertExpectations(t)
		})
	}
}
