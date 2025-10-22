package familyservice

import (
	"cashly/internal/entity"
	"cashly/internal/pkg/sl"
	"cashly/internal/service/familyservice/mocks"
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestCases struct {
	Name      string
	ReturnErr bool
}

func newTestLogger() *sl.MyLogger {
	return sl.New(slog.Default(), nil)
}

func TestCreate(t *testing.T) {
	testCases := []TestCases{
		{"success", false},
		{"failed", true},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			mockService := new(mocks.FamilyServiceIfaceMock)

			expectedFamily := &entity.Family{ID: 1, Name: "Test Family", CreatedBy: 42}
			if tc.ReturnErr {
				mockService.On("Create", mock.Anything, mock.AnythingOfType("*entity.Family")).
					Return(nil, errors.New("db failed"))
			} else {
				mockService.On("Create", mock.Anything, mock.MatchedBy(func(f *entity.Family) bool {
					return f.CreatedBy == 42 && f.Name == "Test Family"
				})).Return(expectedFamily, nil)
			}

			svc := New(mockService, newTestLogger())

			got, err := svc.Create(context.Background(), "Test Family", 42)

			if tc.ReturnErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedFamily, got)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetFamiliesByUserID(t *testing.T) {
	testCases := []TestCases{
		{"success", false},
		{"failed", true},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			mockProvider := new(mocks.FamilyProviderMock)

			expectedFamilies := []entity.Family{{ID: 1, CreatedBy: 42, Name: "Test Family"}}
			if tc.ReturnErr {
				mockProvider.On("GetFamiliesByUserID", mock.Anything, int64(42)).
					Return(nil, errors.New("db failed"))
			} else {
				mockProvider.On("GetFamiliesByUserID", mock.Anything, int64(42)).Return(expectedFamilies, nil)
			}

			svc := &FamilyService{
				familyProvider: mockProvider,
				sl:             newTestLogger(),
			}

			got, err := svc.GetFamiliesByUserID(context.Background(), int64(42))

			if tc.ReturnErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedFamilies, got)
			}
			mockProvider.AssertExpectations(t)
		})
	}
}

func TestGetFamilyByCode(t *testing.T) {
	testCases := []TestCases{
		{"success", false},
		{"failed", true},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			mockProvider := new(mocks.FamilyProviderMock)

			expectedFamily := &entity.Family{ID: 1, CreatedBy: 42, Name: "Test Family"}
			expectedExpiresAt := time.Now().Add(48 * time.Hour)

			if tc.ReturnErr {
				mockProvider.On("GetFamilyByCode", mock.Anything, "ABC123").
					Return(nil, time.Time{}, errors.New("db failed"))
			} else {
				mockProvider.On("GetFamilyByCode", mock.Anything, "ABC123").Return(expectedFamily, expectedExpiresAt, nil)
			}

			svc := &FamilyService{
				familyProvider: mockProvider,
				sl:             newTestLogger(),
			}

			got, gotExpiresAt, err := svc.GetFamilyByCode(context.Background(), "ABC123")

			if tc.ReturnErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedFamily, got)
				assert.Equal(t, expectedExpiresAt, gotExpiresAt)
			}
			mockProvider.AssertExpectations(t)
		})
	}
}

func TestGetFamilyByID(t *testing.T) {
	testCases := []TestCases{
		{"success", false},
		{"failed", true},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			mockProvider := new(mocks.FamilyProviderMock)

			expectedFamily := &entity.Family{ID: 1, CreatedBy: 42, Name: "Test Family"}

			if tc.ReturnErr {
				mockProvider.On("GetFamilyByID", mock.Anything, 1).
					Return(nil, errors.New("db failed"))
			} else {
				mockProvider.On("GetFamilyByID", mock.Anything, 1).Return(expectedFamily, nil)
			}

			svc := &FamilyService{
				familyProvider: mockProvider,
				sl:             newTestLogger(),
			}

			got, err := svc.GetFamilyByID(context.Background(), 1)

			if tc.ReturnErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedFamily, got)
			}
			mockProvider.AssertExpectations(t)
		})
	}
}

func TestGenerateInviteCode(t *testing.T) {
	mockService := new(mocks.FamilyServiceIfaceMock)

	svc := New(mockService, newTestLogger())

	code, err := svc.GenerateInviteCode()

	assert.NoError(t, err)
	assert.Len(t, code, codeLength)
}
