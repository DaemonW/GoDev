package crypto

import "testing"

func TestPKCS7Unpadding(t *testing.T) {
	str := "fwlegljgllfjlgj;leshl;res4983gdkhklz"
	b := PKCS7Padding([]byte(str), 32)
	b=PKCS7Unpadding(b)
	if len(b)!=len([]byte(str)){
		t.Error("length not match")
	}
}

func TestPKCS7Padding(t *testing.T) {
	str := "fwlegljgllfjlgj;leshl;reskerhjslhlsjhlq"
	b := PKCS7Padding([]byte(str), 16)
	b=PKCS7Unpadding(b)
	if len(b)!=len([]byte(str)){
		t.Error("length not match")
	}
}
