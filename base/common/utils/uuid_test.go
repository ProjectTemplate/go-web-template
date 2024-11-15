package utils

import (
	"github.com/google/uuid"
	"testing"
)

func BenchmarkUUID(b *testing.B) {
	//0.0000013 ns/op
	b.Run("UUID", func(b *testing.B) {
		uuid.New().String()
	})

	//0.0000004 ns/op
	b.Run("UUID_V6", func(b *testing.B) {
		uuid.NewV6()
	})

	//0.0000025 ns/op
	b.Run("UUID_V7", func(b *testing.B) {
		uuid.NewV7()
	})
}
