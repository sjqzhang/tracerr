package main

import (
	"io/ioutil"

	"github.com/sjqzhang/tracerr"
)

func main() {
	if err := read(); err != nil {
		tracerr.PrintSourceColor(err)
	}
}

func read() error {
	return readNonExistent()
}

func readNonExistent() error {
	_, err := ioutil.ReadFile("/tmp/non_existent_file")
	// Add stack trace to existing error, no matter if it's nil.
	return tracerr.Wrap(err)
}
