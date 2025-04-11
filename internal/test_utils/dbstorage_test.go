package testutils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

func TestDBStorage(t *testing.T) {
	storage, logger := testutils.SetupDBStorage(t)

	assert.NotNil(t, storage)
	assert.NotNil(t, logger)
}
