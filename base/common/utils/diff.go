package utils

import (
	"github.com/sergi/go-diff/diffmatchpatch"
)

// DiffText 比较两个文本
//
// 相同的部分：输出黑色
//
// 不同的部分：textA的内容输出红色，textB的内容输出绿色
func DiffText(textA, textB string) string {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(textA, textB, false)

	return dmp.DiffPrettyText(diffs)
}
