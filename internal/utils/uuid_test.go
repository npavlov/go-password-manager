package utils_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/npavlov/go-password-manager/internal/utils"
)

func TestGetIdFromString_ValidUUID(t *testing.T) {
	str := uuid.New().String()

	result := utils.GetIdFromString(str)

	assert.True(t, result.Valid)
	assert.Equal(t, str, result.String())
}

func TestGetIdFromString_InvalidUUID(t *testing.T) {
	str := "not-a-uuid"

	result := utils.GetIdFromString(str)

	assert.False(t, result.Valid)
}
