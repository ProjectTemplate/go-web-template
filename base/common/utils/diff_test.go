package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDiffTextEmpty(t *testing.T) {
	text1 := ""
	text2 := ""

	diffText := DiffText(text1, text2)

	assert.Equal(t, text1, diffText)
}

func TestDiffTextSame(t *testing.T) {
	text1 := "Hello A!"
	text2 := "Hello A!"

	diffText := DiffText(text1, text2)

	assert.Equal(t, text1, diffText)
}

func TestDiffTextSomeSame(t *testing.T) {
	textA := "Hello 张三!"
	textB := "Hello Bob"
	diffText := DiffText(textA, textB)

	fmt.Println(diffText)
	assert.NotEqual(t, textA, diffText)
	assert.NotEqual(t, textB, diffText)
}

func TestDiffTextAllDifferent(t *testing.T) {
	textA := "你好啊！"
	textB := "今天吃什么？"
	diffText := DiffText(textA, textB)

	fmt.Println(diffText)
	assert.NotEqual(t, textA, diffText)
	assert.NotEqual(t, textB, diffText)
}
