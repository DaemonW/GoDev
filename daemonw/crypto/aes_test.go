package crypto

import (
	"testing"
)

const (
	ivSeed = "fwelgojgweogj"
)

func genKey() []byte {
	return GetHash([]byte(ivSeed), "SHA256")
}

func genIv() []byte {
	return GetHash([]byte(ivSeed), "MD5")
}

func TestAesCbcEnc(t *testing.T) {
	plainText := "test aes cbc encrypt!"
	key, iv := genKey(), genIv()
	enc, err := AesCbcEnc([]byte(plainText), key, iv)
	if err != nil {
		t.Error("aes encrypt error")
	}
	dec, err := AesCbcDec(enc, key, iv)
	if string(dec) != plainText {
		t.Error("aes decrypt error")
	}
}


func TestAesCbcDec(t *testing.T) {
	plainText := "test aes cbc decrypt!"
	key, iv := genKey(), genIv()
	enc, err := AesCbcEnc([]byte(plainText), key, iv)
	if err != nil {
		t.Error("aes encrypt error")
	}
	dec, err := AesCbcDec(enc, key, iv)
	if string(dec) != plainText {
		t.Error("aes decrypt error")
	}
}




func TestAesCfbEnc(t *testing.T) {
	plainText := "test aes cfb encrypt!"
	key, iv := genKey(), genIv()
	enc, err := AesCfbEnc([]byte(plainText), key, iv)
	if err != nil {
		t.Error("aes encrypt error")
	}
	dec, err := AesCfbDec(enc, key, iv)
	if string(dec) != plainText {
		t.Error("aes decrypt error")
	}
}


func TestAesCfbDec(t *testing.T) {
	plainText := "test aes cfb decrypt!"
	key, iv := genKey(), genIv()
	enc, err := AesCfbEnc([]byte(plainText), key, iv)
	if err != nil {
		t.Error("aes encrypt error")
	}
	dec, err := AesCfbDec(enc, key, iv)
	if string(dec) != plainText {
		t.Error("aes decrypt error")
	}
}


func TestAesCtrEnc(t *testing.T) {
	plainText := "test aes ctr encrypt!"
	key, iv := genKey(), genIv()
	enc, err := AesCtrEnc([]byte(plainText), key, iv)
	if err != nil {
		t.Error("aes encrypt error")
	}
	dec, err := AesCtrDec(enc, key, iv)
	if string(dec) != plainText {
		t.Error("aes decrypt error")
	}
}


func TestAesCtrDec(t *testing.T) {
	plainText := "test aes ctr decrypt!"
	key, iv := genKey(), genIv()
	enc, err := AesCtrEnc([]byte(plainText), key, iv)
	if err != nil {
		t.Error("aes encrypt error")
	}
	dec, err := AesCtrDec(enc, key, iv)
	if string(dec) != plainText {
		t.Error("aes decrypt error")
	}
}
