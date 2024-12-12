package utils

import (
	"github.com/google/go-cmp/cmp"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// DiffText 比较两个文本，当数据相同的时候返回文本内容，当数据不同的时候返回用颜色标记的内容
//
// 相同的部分：输出黑色
//
// 不同的部分：textA的内容输出红色，textB的内容输出绿色
func DiffText(textA, textB string) string {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(textA, textB, false)

	return dmp.DiffPrettyText(diffs)
}

// DiffInterface 比较任意类型数据，当数据相同的时候返回空字符串，当数据不同的时候返回不同的内容
func DiffInterface(dataA, dataB interface{}) string {
	diff := cmp.Diff(dataA, dataB)
	return diff
}
