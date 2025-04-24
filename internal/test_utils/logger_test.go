package testutils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

func TestGetLogger(t *testing.T) {
	t.Parallel()

	logger := testutils.GetTLogger()

	assert.NotNil(t, logger)
}
