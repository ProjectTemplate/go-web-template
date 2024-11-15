package utils

import (
	"strconv"
	"strings"
	"sync/atomic"
)

type Span struct {
	parentSpan string

	number atomic.Int64
}

// NewSpan 创建一个新的 Span，每次处理请求的时候理论上只生成一个 Span 即可
func NewSpan(parentSpan string) *Span {
	span := Span{
		parentSpan: parentSpan,
		number:     atomic.Int64{},
	}
	return &span
}

// IncreaseAndGet 自增并获取 Span，每次获取的 Span 都是唯一的
func (t *Span) IncreaseAndGet() string {
	spanNumber := t.number.Add(1)
	spanStr := strconv.FormatInt(spanNumber, 10)

	if t.parentSpan == "" {
		return spanStr
	}

	return strings.Join([]string{t.parentSpan, spanStr}, ".")
}
