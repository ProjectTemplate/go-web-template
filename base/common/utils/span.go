package utils

import (
	"context"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type Span struct {
	parentSpan string
	name       string
	number     atomic.Int64
	startTime  time.Time
	entTime    time.Time
}

// NewSpan 创建一个新的 Span，每次处理请求的时候理论上只生成一个 Span 即可
func NewSpan(parentSpan, name string) *Span {
	span := &Span{
		parentSpan: parentSpan,
		name:       name,
		number:     atomic.Int64{},
		startTime:  time.Now(),
	}

	//从一开始
	span.number.Store(1)
	return span
}

func (s *Span) Child(childName string) *Span {
	spanNumber := s.number.Load()
	formatInt := strconv.FormatInt(spanNumber, 10)

	return &Span{
		parentSpan: strings.Join([]string{s.parentSpan, formatInt}, "."),
		name:       childName,
		number:     atomic.Int64{},
	}
}

// Next 获取下一个Span
func (s *Span) Next(name string) *Span {
	nextSpanNumber := s.number.Add(1)
	span := Span{
		parentSpan: s.parentSpan,
		name:       name,
		number:     atomic.Int64{},
	}
	span.number.Store(nextSpanNumber)

	return &span
}

func (s *Span) Span() string {
	load := s.number.Load()
	return strings.Join([]string{s.parentSpan, strconv.FormatInt(load, 10)}, ".")
}

func (s *Span) End() {
	s.entTime = time.Now()
}

func Next(ctx context.Context, name string) *Span {
	span, ok := ctx.Value("span").(*Span)
	if ok {
		return span.Next(name)
	}

	return nil
}

func Child(ctx context.Context, name string) *Span {
	span, ok := ctx.Value("span").(*Span)
	if ok {
		return span.Child(name)
	}

	return nil
}
