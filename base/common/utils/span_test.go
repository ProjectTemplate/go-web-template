package utils

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
)

func TestSpan(t *testing.T) {
	span := NewSpan("")

	spanStr := span.IncreaseAndGet()
	assert.Equal(t, "1", spanStr)

	spanStr = span.IncreaseAndGet()
	assert.Equal(t, "2", spanStr)
}

func TestWithParent(t *testing.T) {
	span := NewSpan("1")

	spanStr := span.IncreaseAndGet()
	assert.Equal(t, "1.1", spanStr)

	spanStr = span.IncreaseAndGet()
	assert.Equal(t, "1.2", spanStr)
}

func TestParallel(t *testing.T) {
	count := 1000000
	goRoutineCount := 10

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(goRoutineCount)

	span := NewSpan("1")
	for i := 0; i < goRoutineCount; i++ {
		go invokeIncreaseAndGet(span, count, waitGroup)
	}

	waitGroup.Wait()

	spanResult := strconv.Itoa(goRoutineCount*count + 1)
	assert.Equal(t, "1."+spanResult, span.IncreaseAndGet())
}

func invokeIncreaseAndGet(span *Span, count int, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	for i := 0; i < count; i++ {
		span.IncreaseAndGet()
	}
}
