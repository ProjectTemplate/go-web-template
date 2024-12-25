package utils

import (
	"bufio"
	"os"
)

// FileLineScanner 返回一个文件的行扫描器
func FileLineScanner(filePath string) (*bufio.Scanner, error) {
	openFile, err := os.Open(filePath)
	if err != nil {
		return &bufio.Scanner{}, err
	}

	return bufio.NewScanner(openFile), nil
}
