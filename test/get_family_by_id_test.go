package test

import (
	"cashly/internal/entity"
	familyservice "cashly/internal/service/family"
	"cashly/internal/service/family/mocks"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFamilyByID(t *testing.T) {

	const fid = 1

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
				CreatedBy: 42,
				Name:      "Test Family",
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

			mockService.On("GetByID", context.Background(), fid).
				Return(tt.mockFamily, tt.mockRepoErr)

			svc := familyservice.New(mockService, nil, newTestLogger())

			gotFamily, err := svc.GetByID(context.Background(), fid)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.mockFamily, gotFamily)

			mockService.AssertExpectations(t)
		})
	}
}
