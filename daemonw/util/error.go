package util

import (
	"log"
	"runtime/debug"
)

func FatalIfErr(err error, trace bool) {
	if err != nil {
		log.Fatal(err)
		if trace {
			printTrace()
		}
	}
}

func printTrace() {
	debug.PrintStack()
}

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
