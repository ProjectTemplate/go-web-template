package utils

import (
	"github.com/google/uuid"
)

func UUID() string {
	v7, err := uuid.NewV7()
	if err == nil {
		return v7.String()
	}

	v6, err := uuid.NewV6()
	if err == nil {
		return v6.String()
	}

	return uuid.New().String()
}
