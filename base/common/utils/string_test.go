package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFillZero(t *testing.T) {
	assert.Equal(t, "1", FillZero("1", 0))
	assert.Equal(t, "1", FillZero("1", 1))
	assert.Equal(t, "01", FillZero("1", 2))
	assert.Equal(t, "001", FillZero("1", 3))

	assert.Equal(t, "012", FillZero("12", 3))
	assert.Equal(t, "0012", FillZero("12", 4))
	assert.Equal(t, "00012", FillZero("12", 5))
}
