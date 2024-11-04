package response

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncode(t *testing.T) {
	code := NewCode("001", "002", "003", "004")
	code = code.WithCode("0005")

	assert.Equal(t, "0010020030040005", Encode(code))

	code = code.WithCode("0006")
	assert.Equal(t, "0010020030040006", Encode(code))
}

func TestDecode(t *testing.T) {

	decode := Decode("0010020030040005")

	assert.Equal(t, "001", decode.Company)
	assert.Equal(t, "002", decode.Department)
	assert.Equal(t, "003", decode.Project)
	assert.Equal(t, "004", decode.ProjectModule)
	assert.Equal(t, "0005", decode.Code)
}
