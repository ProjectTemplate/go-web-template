package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStructToMap(t *testing.T) {
	type friends struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	type data struct {
		Name     string            `json:"name"`
		Age      int               `json:"age"`
		Friends  []friends         `json:"friends"`
		Relation map[string]string `json:"relation"`
	}
	d := data{
		Name: "test",
		Age:  18,
		Friends: []friends{
			{Name: "name1", Age: 18},
			{Name: "name2", Age: 19},
		},
		Relation: map[string]string{
			"father": "father",
			"mother": "mother",
		},
	}

	result := StructToMap(d)

	assert.Equal(t, "test", result["name"].(string))
	assert.Equal(t, 18, result["age"].(int))
}
