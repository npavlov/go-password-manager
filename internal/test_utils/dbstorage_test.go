package testutils_test

import (
	"testing"

	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestDBStorage(t *testing.T) {
	storage, logger := testutils.SetupDBStorage(t)

	assert.NotNil(t, storage)
	assert.NotNil(t, logger)
}
