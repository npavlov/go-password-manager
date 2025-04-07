package testutils_test

import (
	"testing"

	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestGetLogger(t *testing.T) {
	logger := testutils.GetTLogger()

	assert.NotNil(t, logger)
}
