package crypto

import (
	"errors"
	"fmt"
)

var (
	hexSeed = [16]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
)

//Int2Bytes get the byte array of the integer data. high bit for high index of the array
//num is a int64 type to accept other integers,such as uint16,etc...
//length mean the data size,such as int64==8, int32==4,uint16==2 etc...
func Int2Bytes(num int64, bitNum int, bigEndian bool) []byte {
	size := bitNum / 8
	data := make([]byte, size)
	if bigEndian {
		for i := size - 1; i >= 0; i-- {
			data[i] = byte(num & 0xff)
			num = num >> 8
		}
	} else {
		for i := 0; i < size; i++ {
			data[i] = byte(num & 0xff)
			num = num >> 8
		}
	}
	return data
}

//Bytes2String convert byte array to hex string
func Bytes2Hex(data []byte) (str string) {
	if data == nil {
		return ""
	}
	size := len(data)
	b := make([]byte, size*2)
	for i := 0; i < size; i++ {
		b[i*2+1] = hexSeed[data[i]&0x0f]
		b[i*2] = hexSeed[data[i]>>4&0x0f]
	}
	return string(b)
}

//Hex2Bytes parse the hex string, convert each two
//hex chars into a byte
func Hex2Bytes(hexStr string) ([]byte, error) {
	var err error
	hex := []byte(hexStr)
	size := len(hex)
	if size%2 != 0 {
		return nil, errors.New("illegal hex string length")
	}
	data := make([]byte, size/2)
	var tmp, al, ah byte
	for i := 0; i < size; i++ {
		if tmp, err = hex2Decimal(hex[i]); err != nil {
			return nil, fmt.Errorf("illegal hex character '%v' at index %d", hex[i], i)
		}
		if i%2 == 1 {
			al = tmp
			data[i/2] = ah<<4 | al
		} else {
			ah = tmp
		}
	}
	return data, nil
}

//hex2Decimal convert a hex char to decimal value
func hex2Decimal(c byte) (uint8, error) {
	switch {
	case c >= 'a' && c <= 'f':
		return c&0x0f + 9, nil
	case c >= 'A' && c <= 'F':
		return c&0x0f + 9, nil
	case c >= '0' && c <= '9':
		return c & 0x0f, nil
	default:
		return 0xff, errors.New("illegal hex character")
	}
}
