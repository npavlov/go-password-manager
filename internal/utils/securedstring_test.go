package utils_test

import (
	"testing"

	"github.com/npavlov/go-password-manager/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewStringAndGet(t *testing.T) {
	s := utils.NewString("secret")

	assert.NotNil(t, s)
	assert.Equal(t, "secret", s.Get())
}

func TestSetValueAndGet(t *testing.T) {
	s := utils.NewString("initial")
	s.Set("updated")

	assert.Equal(t, "updated", s.Get())
}

func TestIsEquals_SameValueSameKey(t *testing.T) {
	a := utils.NewString("match")
	b := utils.NewString("match")

	// Ensure both have the same key
	b.SetKey(a.GetSelf().Key)
	b.Set("match")

	assert.True(t, a.IsEquals(b))
}

func TestIsEquals_DifferentValue(t *testing.T) {
	a := utils.NewString("value1")
	b := utils.NewString("value2")

	b.SetKey(a.GetSelf().Key)

	assert.False(t, a.IsEquals(b))
}

func TestIsEquals_DifferentKey(t *testing.T) {
	a := utils.NewString("data")
	b := utils.NewString("data")

	b.RandomizeKey() // Now keys are different

	assert.True(t, a.IsEquals(b))
}

func TestRandomizeKey_ChangesKey(t *testing.T) {
	s := utils.NewString("rotate-key")

	oldKey := s.GetSelf().Key
	s.RandomizeKey()

	assert.NotEqual(t, oldKey, s.GetSelf().Key)
	assert.Equal(t, "rotate-key", s.Get()) // Ensure value is still correct
}
