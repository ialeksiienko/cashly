package test

import (
	familyservice "cashly/internal/service/family"
	"cashly/internal/service/family/mocks"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var codeLength = 6

func TestGenerateFamilyInviteCode(t *testing.T) {
	mockService := new(mocks.FamilyIfaceMock)

	svc := familyservice.New(mockService, nil, newTestLogger())

	code, err := svc.GenerateInviteCode()

	assert.NoError(t, err)
	t.Log(fmt.Sprintf("code: %s, length: %d", code, len(code)))
	assert.Len(t, code, codeLength)
}
