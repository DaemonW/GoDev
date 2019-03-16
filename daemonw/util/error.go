package util

import (
	dlog "log"
	"runtime/debug"
)

var DEBUG bool = false

func CheckFatal(err error) {
	if err != nil {
		dlog.Fatal(err)

		debug.PrintStack()
	}
}

func PrintStackTrace() {
	debug.PrintStack()
}

func PanicIfErr(err error){
	if err!=nil{
		panic(err)
	}
}
