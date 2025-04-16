package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/npavlov/go-password-manager/internal/utils"
)

func TestNewStringAndGet(t *testing.T) {
	t.Parallel()

	str := utils.NewString("secret")

	assert.NotNil(t, str)
	assert.Equal(t, "secret", str.Get())
}

func TestSetValueAndGet(t *testing.T) {
	t.Parallel()

	str := utils.NewString("initial")
	str.Set("updated")

	assert.Equal(t, "updated", str.Get())
}

func TestIsEquals_SameValueSameKey(t *testing.T) {
	t.Parallel()

	val1 := utils.NewString("match")
	val2 := utils.NewString("match")

	// Ensure both have the same key
	val2.SetKey(val1.GetSelf().Key)
	val2.Set("match")

	assert.True(t, val1.IsEquals(val2))
}

func TestIsEquals_DifferentValue(t *testing.T) {
	t.Parallel()

	a := utils.NewString("value1")
	b := utils.NewString("value2")

	b.SetKey(a.GetSelf().Key)

	assert.False(t, a.IsEquals(b))
}

func TestIsEquals_DifferentKey(t *testing.T) {
	t.Parallel()

	a := utils.NewString("data")
	b := utils.NewString("data")

	b.RandomizeKey() // Now keys are different

	assert.True(t, a.IsEquals(b))
}

func TestRandomizeKey_ChangesKey(t *testing.T) {
	t.Parallel()

	str := utils.NewString("rotate-key")

	oldKey := str.GetSelf().Key
	str.RandomizeKey()

	assert.NotEqual(t, oldKey, str.GetSelf().Key)
	assert.Equal(t, "rotate-key", str.Get()) // Ensure value is still correct
}
