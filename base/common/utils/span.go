package utils

import (
	"strconv"
	"strings"
	"sync/atomic"
)

type Span struct {
	parentSpan string
	name       string
	number     atomic.Int64
}

// NewSpan 创建一个新的 Span，每次处理请求的时候理论上只生成一个 Span 即可
func NewSpan(parentSpan, name string) *Span {
	span := Span{
		parentSpan: parentSpan,
		name:       name,
		number:     atomic.Int64{},
	}
	return &span
}

func (s *Span) Get() string {
	spanNumber := s.number.Load()
	spanStr := strconv.FormatInt(spanNumber, 10)
	if s.parentSpan == "" {
		return spanStr
	}

	return strings.Join([]string{s.parentSpan, spanStr}, ".")
}

// IncreaseAndGet 自增并获取 Span，每次获取的 Span 都是唯一的
func (s *Span) IncreaseAndGet() string {
	spanNumber := s.number.Add(1)
	spanStr := strconv.FormatInt(spanNumber, 10)

	if s.parentSpan == "" {
		return spanStr
	}

	return strings.Join([]string{s.parentSpan, spanStr}, ".")
}
