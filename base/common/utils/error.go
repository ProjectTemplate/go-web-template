package utils

import "fmt"

func PanicAndPrintIfNotNil(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
