package utils

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		go invokeChild(span, count, waitGroup)
	}

	waitGroup.Wait()

	assert.Equal(t, int64(count*goRoutineCount), span.childSpanNumber.Load())
}

func invokeChild(span *Span, count int, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	for i := 0; i < count; i++ {
		span.Child("child_" + strconv.Itoa(i))
	}
}
