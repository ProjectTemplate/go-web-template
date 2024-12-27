package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	test(t)
}

func test(t *testing.T) {
	stack := CallerStack(0, 0)

	assert.NotNil(t, stack)
	//fmt.Println(stack.String())
}

func TestGetParentCallerMethodName(t *testing.T) {
	callerMethodName := GetParentCallerMethodName()
	assert.Equal(t, "TestGetParentCallerMethodName", callerMethodName)

	innerFunc(t)
}

func innerFunc(t *testing.T) {
	callerMethodName := GetParentCallerMethodName()
	assert.Equal(t, "innerFunc", callerMethodName)
}

// BenchmarkCallerStack 1593 ns/op
// Warning 注意性能问题，不要滥用
func BenchmarkGetParentCallerMethodName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetParentCallerMethodName()
	}
}

// BenchmarkEmpty 7.641 ns/op
func BenchmarkEmpty(b *testing.B) {
	data := make(map[string]string)
	for i := 0; i < b.N; i++ {
		data["a"] = "a"
	}
}
