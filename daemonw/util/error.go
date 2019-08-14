package util

import (
	"runtime/debug"
)

func PrintTrace() {
	debug.PrintStack()
}

func StackInfo() string{
	return string(debug.Stack())
}

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
