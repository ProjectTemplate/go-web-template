package utils

import (
	"bytes"
)

// FillZeroToNumberString 填充0到字符串前面，直到字符串长度达到totalLength
func FillZeroToNumberString(str string, totalLength int) string {
	if len(str) >= totalLength {
		return str
	}

	buffer := bytes.NewBuffer(make([]byte, 0, totalLength))
	for i := len(str) + 1; i <= totalLength; i++ {
		buffer.WriteString("0")
	}
	buffer.WriteString(str)

	return buffer.String()
}
