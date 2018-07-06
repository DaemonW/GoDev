#!/bin/bash
case $1 in

"linux32")
export GOOS="linux"
CGO_ENABLE=0
GOARCH="in386"
;;
"linux64")
export GOOS="linux"
CGO_ENABLE=0
GOARCH="amd64"
;;
"win32")
export GOOS="windows"
CGO_ENABLE=0
GOARCH="in386"
;;
"win64")
export GOOS="windows"
CGO_ENABLE=0
GOARCH="amd64"
;;
"")
;;
*)
echo "unknown archtecture"
;;
esac
