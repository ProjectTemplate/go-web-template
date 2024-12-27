package kafka

import (
	"testing"

	"github.com/segmentio/kafka-go"
)

var partitions = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

func BenchmarkKeyBalancer(b *testing.B) {
	message := kafka.Message{}
	message.Key = []byte("123456789")
	message.Value = []byte("123456789")

	//0.0000015 ns/op
	hash := &kafka.Hash{}
	b.Run("HashBalancer", func(b *testing.B) {
		runBalancer(hash, message)
	})

	//0.0000005 ns/op
	crc32Balancer := &kafka.CRC32Balancer{}
	b.Run("crc32Balancer", func(b *testing.B) {
		runBalancer(crc32Balancer, message)
	})

	//0.0000002 ns/op
	murmur2Balancer := &kafka.Murmur2Balancer{}
	b.Run("murmur2Balancer", func(b *testing.B) {
		runBalancer(murmur2Balancer, message)
	})

	//0.0000014 ns/op
	referenceHash := &kafka.ReferenceHash{}
	b.Run("referenceHash", func(b *testing.B) {
		runBalancer(referenceHash, message)
	})

	//0.0000000 ns/op
	roundRobin := &kafka.RoundRobin{}
	b.Run("roundRobin-random", func(b *testing.B) {
		runBalancer(roundRobin, message)
	})

	//最少消息hash
	//0.0000003 ns/op
	leastBytes := &kafka.LeastBytes{}
	b.Run("leastBytes-random", func(b *testing.B) {
		runBalancer(leastBytes, message)
	})

	//空白对照组
	//0.0000000 ns/op
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
