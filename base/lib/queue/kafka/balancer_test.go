package kafka

import (
	"github.com/segmentio/kafka-go"
	"testing"
)

var partitions = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

func BenchmarkKeyBalancer(b *testing.B) {
	message := kafka.Message{}
	message.Key = []byte("123456789")
	message.Value = []byte("123456789")

	hash := &kafka.Hash{}
	b.Run("HashBalancer", func(b *testing.B) {
		runBalancer(hash, message)
	})

	crc32Balancer := &kafka.CRC32Balancer{}
	b.Run("crc32Balancer", func(b *testing.B) {
		runBalancer(crc32Balancer, message)
	})

	murmur2Balancer := &kafka.Murmur2Balancer{}
	b.Run("murmur2Balancer", func(b *testing.B) {
		runBalancer(murmur2Balancer, message)
	})

	referenceHash := &kafka.ReferenceHash{}
	b.Run("referenceHash", func(b *testing.B) {
		runBalancer(referenceHash, message)
	})

	roundRobin := &kafka.RoundRobin{}
	b.Run("roundRobin-random", func(b *testing.B) {
		runBalancer(roundRobin, message)
	})

	//最少消息hash
	leastBytes := &kafka.LeastBytes{}
	b.Run("leastBytes-random", func(b *testing.B) {
		runBalancer(leastBytes, message)
	})

	b.Run("do-nothing", func(b *testing.B) {
		// do nothing
	})
}

func runBalancer(balancer kafka.Balancer, message kafka.Message) {
	balance := balancer.Balance(message, partitions...)

	if balance < 0 || balance > 11 {
		panic("balance error")
	}
}
