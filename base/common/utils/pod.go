package utils

import (
	"github.com/google/uuid"
	"os"
)

func GetHostName() string {
	name, err := os.Hostname()
	if err != nil {
		return "hostname-" + uuid.New().String()
	}

	return name
}
