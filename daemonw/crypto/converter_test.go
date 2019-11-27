package crypto

import (
	"testing"
)

func TestBytes2Hex(t *testing.T) {
	b := []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef}
	str := "1234567890abcdef"
	s := Bytes2Hex(b)
	if s != str {
		t.Error("hex string not match")
	}
}

func TestHex2Bytes(t *testing.T) {
	b := []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef}
	str := "1234567890abcdef"
	errStr := "1234567890abcdef0"
	d, err := Hex2Bytes(str)
	if err != nil {
		t.Error(err)
	}
	if len(d) != len(b) {
		t.Error("length not match")
	}

	for i := 0; i < len(b); i++ {
		if b[i] != d[i] {
			t.Error("byte not match")
			return
		}
	}
	b, err = Hex2Bytes(errStr)
	if err==nil{
		t.Error("parse hex string error")
	}
}

func TestInt2Bytes(t *testing.T) {
	num := 0x12345678
	littleE := []byte{0x78, 0x56, 0x34, 0x12}
	bigE := []byte{0x12, 0x34, 0x56, 0x78}
	b := Int2Bytes(int64(num), 32, true)
	for i := 0; i < 4; i++ {
		if b[i] != bigE[i] {
			t.Error("bit endian not match")
		}
	}
	b = Int2Bytes(int64(num), 32, false)
	for i := 0; i < 4; i++ {
		if b[i] != littleE[i] {
			t.Error("little endian not match")
		}
	}
}
