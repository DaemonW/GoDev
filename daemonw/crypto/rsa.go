package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
)

func GenRsaKeyPair(keySize int) (pub []byte, pri []byte, err error) {
	key, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, nil, err
	}
	pri, err = x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, nil, err
	}
	pubKey := key.Public()
	pub, err = x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, nil, err
	}
	return pub, pri, nil
}

func ParsePublicKey(d []byte) (*rsa.PublicKey, error) {
	pub, err := x509.ParsePKIXPublicKey(d)
	if err != nil {
		return nil, err
	}
	if _, ok := pub.(*rsa.PublicKey); !ok {
		return nil, errors.New("illegal rsa public")
	} else {
		return pub.(*rsa.PublicKey), err
	}
}

func ParsePrivateKey(d []byte) (*rsa.PrivateKey, error) {
	pub, err := x509.ParsePKCS8PrivateKey(d)
	if err != nil {
		return nil, err
	}
	if _, ok := pub.(*rsa.PrivateKey); !ok {
		return nil, errors.New("illegal rsa private")
	} else {
		return pub.(*rsa.PrivateKey), err
	}
}

func RsaEnc(key *rsa.PublicKey, src []byte) ([]byte, error) {
	r, err := rsa.EncryptPKCS1v15(rand.Reader, key, src)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func RsaDec(key *rsa.PrivateKey, src []byte) ([]byte, error) {
	r, err := rsa.DecryptPKCS1v15(rand.Reader, key, src)
	if err != nil {
		return nil, err
	}
	return r, nil
}
