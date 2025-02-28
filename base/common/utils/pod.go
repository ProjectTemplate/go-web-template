package utils

import (
	"os"

	"github.com/google/uuid"
)

func GetHostName() string {
	name, err := os.Hostname()
	if err != nil {
		return "hostname-" + uuid.New().String()
	}

	return name
}
