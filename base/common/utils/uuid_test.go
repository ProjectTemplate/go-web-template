package utils

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkUUID(b *testing.B) {
	//0.0000013 ns/op
	b.Run("UUID", func(b *testing.B) {
		id := uuid.New().String()
		assert.NotNil(b, id)
	})

	//0.0000004 ns/op
	b.Run("UUID_V6", func(b *testing.B) {
		id, err := uuid.NewV6()
		assert.Nil(b, err)
		assert.NotNil(b, id)
	})

	//0.0000025 ns/op
	b.Run("UUID_V7", func(b *testing.B) {
		id, err := uuid.NewV7()
		assert.Nil(b, err)
		assert.NotNil(b, id)
	})
}
