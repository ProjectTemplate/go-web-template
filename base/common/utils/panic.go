package utils

import (
	"fmt"
)

func RecoverWithFmt() {
	if err := recover(); err != nil {
		fmt.Printf("PanicRecoverWithFmt %#v", err)
	}
}
