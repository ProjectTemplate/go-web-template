package utils

import (
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type Span struct {
	parentSpan      string
	name            string
	span            int64
	childSpanNumber atomic.Int64
	startTime       time.Time
	entTime         time.Time
}

// NewSpan 创建一个新的 Span
func NewSpan(parentSpan, name string) *Span {
	span := initSpan(parentSpan, name)
	return span
}

func (s *Span) Child(childName string) *Span {
	nextChildSpanNumber := s.childSpanNumber.Add(1)

	parentSpan := strconv.FormatInt(s.span, 10)
	if s.parentSpan != "" {
		parentSpan = strings.Join([]string{s.parentSpan, parentSpan}, ".")
	}

	span := initSpan(parentSpan, childName)
	span.span = nextChildSpanNumber

	return span
}

func initSpan(parentSpan string, name string) *Span {
	span := Span{
		parentSpan:      parentSpan,
		name:            name,
		span:            1,
		childSpanNumber: atomic.Int64{},
		startTime:       time.Now(),
		entTime:         time.Now(),
	}
	span.childSpanNumber.Store(0)
	return &span
}

func (s *Span) End() {
	s.entTime = time.Now()
}

func (s *Span) Span() string {
	if s.parentSpan == "" {
		return strconv.FormatInt(s.span, 10)
	}

	return strings.Join([]string{s.parentSpan, strconv.FormatInt(s.span, 10)}, ".")
}

func (s *Span) GetParentSpan() string {
	return s.parentSpan
}

func (s *Span) GetName() string {
	return s.name
}

func (s *Span) GetStartTime() int64 {
	return s.startTime.UnixMicro()
}

func (s *Span) GetEndTime() int64 {
	return s.entTime.UnixMicro()
}

func (s *Span) GetDuration() int64 {
	return s.entTime.Sub(s.startTime).Microseconds()
}
