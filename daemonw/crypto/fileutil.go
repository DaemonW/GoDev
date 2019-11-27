package crypto

import (
	"crypto/cipher"
	"errors"
	"io"
	"os"
	"strings"
)

func AesStream(src, dst *os.File, key, iv []byte, mode string, isEncrypt bool) bool {
	rBuff := make([]byte, 8192)
	cBuff := make([]byte, 8192)
	nRead := 0
	var rErr error
	cipherStream, err := getAesSteamCipher(key, iv, mode, isEncrypt)
	if err != nil {
		return false
	}
	for {
		nRead, rErr = src.Read(rBuff)
		if rErr != nil && rErr != io.EOF {
			return false
		}
		cipherStream.XORKeyStream(cBuff[:nRead], rBuff[:nRead])
		_, err = dst.Write(cBuff[:nRead])
		if err != nil {
			return false
		}
		if rErr == io.EOF {
			dst.Sync()
			break
		}
	}
	return true
}

func AesBlock(src, dst *os.File, key, iv []byte, mode string, isEncrypt bool) bool {
	rBuff := make([]byte, 8192)
	cBuff := make([]byte, 8192)
	nRead := 0
	var rErr error
	cipherBlock, err := getAesBlockCipher(key, iv, mode, isEncrypt)
	if err != nil {
		return false
	}
	fi, _ := src.Stat()
	fileLen := fi.Size()
	nLeft := fileLen
	for {
		nRead, rErr = src.Read(rBuff)
		if rErr != nil && rErr != io.EOF {
			return false
		}
		nLeft -= int64(nRead)
		if nLeft == 0 {
			var tmpR []byte
			var tmpW []byte
			if isEncrypt {
				tmpR = PKCS7Padding(rBuff[:nRead], 16)
				tmpW = make([]byte, len(tmpR))
				cipherBlock.CryptBlocks(tmpW, tmpR)
				_, err = dst.Write(tmpW)
				if err != nil {
					return false
				}
			} else {
				cipherBlock.CryptBlocks(cBuff[:nRead], rBuff[:nRead])
				tmpW = PKCS7Unpadding(cBuff[:nRead])
				_, err = dst.Write(tmpW)
				if err != nil {
					return false
				}
			}
		} else {
			cipherBlock.CryptBlocks(cBuff[:nRead], rBuff[:nRead])
			_, err = dst.Write(cBuff[:nRead])
			if err != nil {
				return false
			}
		}

		if rErr == io.EOF || nLeft == 0 {
			dst.Sync()
			break
		}
	}
	return true
}

func aesStreamFile(originPath, newPath string, key, iv []byte, mode string, isEncrypt bool) bool {
	fr, err := os.Open(originPath)
	if err != nil {
		return false
	}
	defer fr.Close()
	fw, err := os.OpenFile(newPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return false
	}
	defer fw.Close()
	return AesStream(fr, fw, key, iv, mode, isEncrypt)
}

//aesBlockFile aes encrypt with block mode, use pkcs7padding
func aesBlockFile(originPath, newPath string, key, iv []byte, mode string, isEncrypt bool) bool {
	fr, err := os.Open(originPath)
	if err != nil {
		return false
	}
	defer fr.Close()
	fw, err := os.OpenFile(newPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return false
	}
	defer fw.Close()
	return AesBlock(fr, fw, key, iv, mode, isEncrypt)
}

func AesCtrEncFile(originPath, newPath string, key, iv []byte) bool {
	return aesStreamFile(originPath, newPath, key, iv, "CTR", true)
}

func AesCtrDecFile(originPath, newPath string, key, iv []byte) bool {
	return aesStreamFile(originPath, newPath, key, iv, "CTR", false)
}

func AesCfbEncFile(originPath, newPath string, key, iv []byte) bool {
	return aesStreamFile(originPath, newPath, key, iv, "CFB", true)
}

func AesCfbDecFile(originPath, newPath string, key, iv []byte) bool {
	return aesStreamFile(originPath, newPath, key, iv, "CFB", false)
}

func AesCbcEncFile(originPath, newPath string, key, iv []byte) bool {
	return aesBlockFile(originPath, newPath, key, iv, "CBC", true)
}

func AesCbcDecFile(originPath, newPath string, key, iv []byte) bool {
	return aesBlockFile(originPath, newPath, key, iv, "CBC", false)
}

func getAesSteamCipher(key, iv []byte, mode string, isEncrypt bool) (cipher.Stream, error) {
	mode = strings.ToUpper(mode)
	switch mode {
	case "CFB":
		return CfbCipher(key, iv, isEncrypt)
	case "CTR":
		return CtrCipher(key, iv)
	}
	return nil, errors.New("unsupported mode")
}

func getAesBlockCipher(key, iv []byte, mode string, isEncrypt bool) (cipher.BlockMode, error) {
	mode = strings.ToUpper(mode)
	switch mode {
	case "CBC":
		return CbcCipher(key, iv, isEncrypt)
	}
	return nil, errors.New("unsupported mode")
}
