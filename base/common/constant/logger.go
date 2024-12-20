package constant

const (
	LoggerKeyType           = "type"
	LoggerKeyParentSpan     = "parentSpan"
	LoggerKeySpan           = "span"
	LoggerKeyName           = "name"
	LoggerKeySpanStartUs    = "spanStartUs"
	LoggerKeySpanEndUs      = "spanEndUs"
	LoggerKeySpanDurationUs = "spanDurationUs"

	LoggerKeyTimestampUs = "timestampUs"
	// LoggerKeyDurationUs 时间间隔，单位微秒
	LoggerKeyDurationUs = "durationUs"

	LoggerKeyRemoteIp    = ContextKeyRemoteIp
	LoggerKeyURL         = ContextKeyURL
	LoggerKeyHost        = ContextKeyHost
	LoggerKeyPath        = ContextKeyPath
	LoggerKeyQuery       = ContextKeyQuery
	LoggerKeyPostForm    = ContextKeyPostForm
	LoggerKeyRequestBody = ContextKeyRequestBody

	LoggerKeyTraceId = HeaderKeyTraceId
)

const (
	LoggerTypeSpan = "span"
	LoggerTypeLog  = "log"
)
