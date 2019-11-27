package crypto

import (
	"testing"
)

var (
	TestData   = "abcdefghijklmnopqrstuvwxyz"
	TestBase64 = "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXo="
)

func TestBase64Encode(t *testing.T) {

}

func TestBase64Encode2Str(t *testing.T) {
	str := Base64Encode2Str([]byte(TestData))
	if str != TestBase64 {
		t.Error("base64 encode result not match")
	}
}

func TestBase64Decode(t *testing.T) {
	d,err:=Base64DecodeStr(TestBase64)
	if err!=nil{
		t.Error(err)
	}
	if string(d)!=TestData{
		t.Error("base64 decode result not match")
	}
}

func TestBase64DecodeStr(t *testing.T) {

}
