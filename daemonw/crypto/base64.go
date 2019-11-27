package crypto

import (
	"encoding/base64"
)

func Base64Encode2Str(d []byte) string {
	return base64.StdEncoding.EncodeToString(d)
}

func Base64Encode(d []byte) []byte {
	enc := make([]byte, base64.StdEncoding.EncodedLen(len(d)))
	base64.StdEncoding.Encode(enc, d)
	return enc
}

func Base64DecodeStr(d string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(d)
}

func Base64Decode(d []byte) ([]byte, error) {
	dec := make([]byte, base64.StdEncoding.DecodedLen(len(d)))
	_, err := base64.StdEncoding.Decode(dec, d)
	return dec, err
}
