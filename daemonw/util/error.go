package util

import (
	"log"
	"runtime/debug"
)

var DEBUG bool = false

func CheckFatal(err error) {
	if err != nil {
		log.Fatal(err)

		debug.PrintStack()
	}
}

func PrintStackTrace() {
	debug.PrintStack()
}

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
