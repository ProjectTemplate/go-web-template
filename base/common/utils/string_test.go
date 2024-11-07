package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFillZero(t *testing.T) {
	assert.Equal(t, "1", FillZeroToNumberString("1", 0))
	assert.Equal(t, "1", FillZeroToNumberString("1", 1))
	assert.Equal(t, "01", FillZeroToNumberString("1", 2))
	assert.Equal(t, "001", FillZeroToNumberString("1", 3))

	assert.Equal(t, "012", FillZeroToNumberString("12", 3))
	assert.Equal(t, "0012", FillZeroToNumberString("12", 4))
	assert.Equal(t, "00012", FillZeroToNumberString("12", 5))
}
