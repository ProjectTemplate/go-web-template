package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStack(t *testing.T) {
	test(t)
}

func test(t *testing.T) {
	stack := CallerStack(0, 0)

	assert.NotNil(t, stack)
	fmt.Println(stack.String())
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
