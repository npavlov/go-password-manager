package utils_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/npavlov/go-password-manager/internal/utils"
)

func TestGetIdFromString_ValidUUID(t *testing.T) {
	t.Parallel()

	str := uuid.New().String()

	result := utils.GetIDFromString(str)

	assert.True(t, result.Valid)
	assert.Equal(t, str, result.String())
}

func TestGetIdFromString_InvalidUUID(t *testing.T) {
	t.Parallel()

	str := "not-a-uuid"

	result := utils.GetIDFromString(str)

	assert.False(t, result.Valid)
}
