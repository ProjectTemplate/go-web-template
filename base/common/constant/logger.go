package constant

const (
	LoggerKeyType           = "type"
	LoggerKeyParentSpan     = "parentSpan"
	LoggerKeySpan           = "span"
	LoggerKeyName           = "name"
	LoggerKeySpanStartUs    = "spanStartUs"
	LoggerKeySpanEndUs      = "spanEndUs"
	LoggerKeySpanDurationUs = "spanDurationUs"

	LoggerKeyTimestamp = "timestamp"
	// LoggerKeyDurationUs 时间间隔，单位微秒
	LoggerKeyDurationUs = "durationUs"

	LoggerKeyTraceId = HeaderKeyTraceId
)

const (
	LoggerTypeSpan = "span"
	LoggerTypeLog  = "log"
)
