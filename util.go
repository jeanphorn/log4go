package log4go

import (
	"fmt"
	"os"
)

func recoverPanic() {
	if e := recover(); e != nil {
		fmt.Printf("Panicing %s\n", e)
	}
}

func fileExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return os.IsExist(err)
	}
	return true
}
