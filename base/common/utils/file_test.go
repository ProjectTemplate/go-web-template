package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestFileLineScanner(t *testing.T) {
	dir := os.TempDir()
	filePath := filepath.Join(dir, "data.txt")
	file, err := os.Create(filePath)
	assert.Nil(t, err)

	dataArray := []string{"Hello", "World", "!"}

	for _, data := range dataArray {
		_, err = file.WriteString(data)
		assert.Nil(t, err)

		_, err = file.WriteString("\n")
		assert.Nil(t, err)
	}

	file.Close()

	scanner, err := FileLineScanner(filePath)
	assert.Nil(t, err)

	result := make([]string, 0, 3)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	assert.Equal(t, dataArray, result)
}
