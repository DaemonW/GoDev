package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

func AesCfbEnc(src, key, iv []byte) ([]byte, error) {
	var data []byte
	cfb, err := CfbCipher(key, iv, true)
	if err != nil {
		return nil, err
	}
	data = make([]byte, len(src))
	cfb.XORKeyStream(data, src)
	return data, nil
}

func AesCfbDec(src, key, iv []byte) ([]byte, error) {
	var data []byte
	cfb, err := CfbCipher(key, iv, false)
	if err != nil {
		return nil, err
	}
	data = make([]byte, len(src))
	cfb.XORKeyStream(data, src)
	return data, nil
}

func AesCtrEnc(src, key, iv []byte) ([]byte, error) {
	var data []byte
	ctr, err := CtrCipher(key, iv)
	if err != nil {
		return nil, err
	}
	data = make([]byte, len(src))
	ctr.XORKeyStream(data, src)
	return data, nil
}

func AesCtrDec(src, key, iv []byte) ([]byte, error) {
	var data []byte
	ctr, err := CtrCipher(key, iv)
	if err != nil {
		return nil, err
	}
	data = make([]byte, len(src))
	ctr.XORKeyStream(data, src)
	return data, nil
}

func AesCbcEnc(src, key, iv []byte) ([]byte, error) {
	var data []byte
	cbc, err := CbcCipher(key, iv, true)
	if err != nil {
		return nil, err
	}
	src = PKCS7Padding(src, cbc.BlockSize()) //use PKCS7 padding
	data = make([]byte, len(src))
	cbc.CryptBlocks(data, src)
	return data, nil
}

func AesCbcDec(src, key, iv []byte) ([]byte, error) {
	var data []byte
	cbc, err := CbcCipher(key, iv, false)
	if err != nil {
		return nil, err
	}
	data = make([]byte, len(src))
	cbc.CryptBlocks(data, src)
	data = PKCS7Unpadding(data) //use PKCS7 unPadding
	return data, nil
}

func CtrCipher(key, iv []byte) (cipher.Stream, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewCTR(block, iv), nil
}

func CfbCipher(key, iv []byte, isEnc bool) (cipher.Stream, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if isEnc {
		return cipher.NewCFBEncrypter(block, iv), nil
	}
	return cipher.NewCFBDecrypter(block, iv), nil
}

func CbcCipher(key, iv []byte, isEnc bool) (cipher.BlockMode, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if isEnc {
		return cipher.NewCBCEncrypter(block, iv), nil
	}
	return cipher.NewCBCDecrypter(block, iv), nil
}
