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

func TestDiffInterfaceString(t *testing.T) {
	textA := "Hello"
	textB := "Hello"
	diffText := DiffInterface(textA, textB)

	assert.Equal(t, "", diffText)
}

func TestDiffInterfaceStructSame(t *testing.T) {
	type person struct {
		Name    string            `json:"name"`
		Age     int               `json:"age"`
		Friends []string          `json:"friends"`
		Extra   map[string]string `json:"extra"`
	}

	dataA := person{
		Name:    "张三",
		Age:     18,
		Friends: []string{"A", "B"},
		Extra:   make(map[string]string),
	}
	dataA.Extra["A"] = "A"
	dataA.Extra["B"] = "B"

	dataB := person{
		Name:    "张三",
		Age:     18,
		Friends: []string{"A", "B"},
		Extra:   make(map[string]string),
	}
	dataB.Extra["A"] = "A"
	dataB.Extra["B"] = "B"

	diffText := DiffInterface(dataA, dataB)

	assert.Equal(t, "", diffText)
}

func TestDiffInterfaceStructDifferent(t *testing.T) {
	type person struct {
		Name    string            `json:"name"`
		Age     int               `json:"age"`
		Friends []string          `json:"friends"`
		Extra   map[string]string `json:"extra"`
	}

	dataA := person{
		Name:    "张三A",
		Age:     15,
		Friends: []string{"A", "B", "C", "E"},
		Extra:   make(map[string]string),
	}
	dataA.Extra["A"] = "A"
	dataA.Extra["B"] = "C"
	dataA.Extra["C"] = "D"
	dataA.Extra["E"] = "F"

	dataB := person{
		Name:    "张三",
		Age:     18,
		Friends: []string{"A", "B", "D", "E"},
		Extra:   make(map[string]string),
	}
	dataB.Extra["A"] = "A"
	dataB.Extra["B"] = "B"
	dataB.Extra["C"] = "C"
	dataB.Extra["H"] = "H"

	diffText := DiffInterface(dataA, dataB)

	fmt.Println(diffText)

	assert.NotEqual(t, "", diffText)
}

func ExampleDiffText() {
	textA := "Hello 张三!"
	textB := "Hello Bob"
	diffText := DiffText(textA, textB)

	fmt.Println(diffText)

	//output:
	//Hello [31m张三![0m[32mBob[0m
}

func ExampleDiffInterface() {
	type person struct {
		Name    string            `json:"name"`
		Age     int               `json:"age"`
		Friends []string          `json:"friends"`
		Extra   map[string]string `json:"extra"`
	}

	dataA := person{
		Name:    "张三A",
		Age:     15,
		Friends: []string{"A", "B", "C", "E"},
		Extra:   make(map[string]string),
	}
	dataA.Extra["A"] = "A"
	dataA.Extra["B"] = "C"
	dataA.Extra["C"] = "D"
	dataA.Extra["E"] = "F"

	dataB := person{
		Name:    "张三",
		Age:     18,
		Friends: []string{"A", "B", "D", "E"},
		Extra:   make(map[string]string),
	}
	dataB.Extra["A"] = "A"
	dataB.Extra["B"] = "B"
	dataB.Extra["C"] = "C"
	dataB.Extra["H"] = "H"

	diffText := DiffInterface(dataA, dataB)

	fmt.Println(diffText)

	//output:
	//utils.person{
	//- 	Name: "张三A",
	//+ 	Name: "张三",
	//- 	Age:  15,
	//+ 	Age:  18,
	//	Friends: []string{
	//		"A",
	//		"B",
	//- 		"C",
	//+ 		"D",
	//		"E",
	//	},
	//	Extra: map[string]string{
	//		"A": "A",
	//- 		"B": "C",
	//+ 		"B": "B",
	//- 		"C": "D",
	//+ 		"C": "C",
	//- 		"E": "F",
	//+ 		"H": "H",
	//	},
	//}
}
