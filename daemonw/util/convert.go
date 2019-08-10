package util

import (
	"errors"
)

var (
	HexSeed = [10]byte{0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39}
)

func Bytes2HexStr(data []byte) string {
	dataLen := len(data)
	str := make([]byte, dataLen*2)
	for i := 0; i < dataLen; i++ {
		str[i*2] = HexSeed[data[i]>>4&0x0f]
		str[i*2+1] = HexSeed[data[i]&0x0f]
	}
	return string(str[:])
}

func HexStr2Bytes(str string) ([]byte, error) {
	d := []byte(str)
	if len(d)%2 != 0 {
		return nil, errors.New("hex characters' length is illegal")
	}
	bytes := make([]byte, len(d)/2)
	var bh, bl byte
	var err error
	for i := 0; i < len(bytes); i++ {
		bh, err = hex2Byte(d[i*2])
		if err != nil {
			return nil, err
		}
		bl, err = hex2Byte(d[i*2+1])
		if err != nil {
			return nil, err
		}
		bytes[i] = bh<<4 | bl
	}
	return bytes, nil
}

func hex2Byte(ch byte) (byte, error) {
	if ch >= 0x30 && ch <= 0x39 {
		return ch & 0x0f, nil;
	}
	if ch >= 0x41 && ch <= 0x5a {
		return ch&0x0f + 9, nil
	}
	if ch >= 0x61 && ch <= 0x7a {
		return ch&0x0f + 9, nil
	}
	return 0, errors.New("illegal hex character")
}
