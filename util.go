package log4go

import "fmt"

func recoverPanic() {
	if e := recover(); e != nil {
		fmt.Printf("Panicing %s\n", e)
	}
}
