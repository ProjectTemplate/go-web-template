package utils

// EmptySpan 空 Span，当从 context 中获取 Span 失败时返回，
// 调用该实例的所有方法都返回空数据，避免链路日志数据错乱
var EmptySpan = initSpan("", "")
