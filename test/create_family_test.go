package test

import (
	"cashly/internal/entity"
	familyservice "cashly/internal/service/family"
	"cashly/internal/service/family/mocks"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateFamily(t *testing.T) {
	tests := []struct {
		name        string
		mockRepoErr error
		mockFamily  *entity.Family
		wantErr     bool
	}{
		{
			name:        "success",
			mockRepoErr: nil,
			mockFamily: &entity.Family{
				ID:        1,
				Name:      "Test Family",
				CreatedBy: 42,
			},
			wantErr: false,
		},
		{
			name:        "database error",
			mockRepoErr: assert.AnError,
			mockFamily:  nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.FamilyIfaceMock)

			mockService.On("Create", context.Background(), mock.AnythingOfType("*entity.Family")).
				Return(tt.mockFamily, tt.mockRepoErr)

			svc := familyservice.New(mockService, nil, newTestLogger())

			gotFamily, err := svc.Create(context.Background(), "Test Family", 42)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.mockFamily, gotFamily)

			mockService.AssertExpectations(t)
		})
	}
}
