package main

import (
	"fmt"
	"github.com/sjqzhang/tracerr"
)

func main() {
	if err := foo(); err != nil {
		//tracerr.PrintSourceColor(err)
		fmt.Println(tracerr.Sprint(err))
	}
}

func foo() error {
	return bar(0)
}

func bar(i int) error {
	if i >= 2 {
		// Create new error with stack trace.
		return tracerr.Errorf("i = %d", i)
	}
	return bar(i + 1)
}
