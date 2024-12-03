package utils

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
)

func TestSpan(t *testing.T) {
	span := NewSpan("", "root")
	assert.Equal(t, "1", span.Span())

	span = span.Next("root")
	assert.Equal(t, "2", span.Span())

	span = span.Next("root")
	assert.Equal(t, "3", span.Span())
}

func TestWithParent(t *testing.T) {
	span := NewSpan("1", "child")
	assert.Equal(t, "1.1", span.Span())

	span = span.Next("child")
	assert.Equal(t, "1.2", span.Span())

	span = span.Next("child")
	assert.Equal(t, "1.3", span.Span())
}

func TestChild(t *testing.T) {
	parent := NewSpan("", "parent")

	child := parent.Child("child")
	assert.Equal(t, "1.1", child.Span())

	grandson := child.Child("grandson")
	assert.Equal(t, "1.1.1", grandson.Span())
}

func TestParallel(t *testing.T) {
	count := 1000000
	goRoutineCount := 10

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(goRoutineCount)

	span := NewSpan("1", "")
	for i := 0; i < goRoutineCount; i++ {
		go invokeNext(span, count, waitGroup)
	}

	waitGroup.Wait()

	spanResult := strconv.Itoa(goRoutineCount*count + 1)

	assert.Equal(t, "1."+spanResult, span.Span())
}

func invokeNext(span *Span, count int, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	for i := 0; i < count; i++ {
		span.Next("")
	}
}
