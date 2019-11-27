package crypto

import "testing"

func TestAesCtrDecFile(t *testing.T) {
	var key [32]byte
	for i := 0; i < 32; i++ {
		key[i] = byte(i)
	}
	var iv [16]byte
	for i := 0; i < 16; i++ {
		iv[i] = byte(i)
	}
	AesCbcEncFile("/home/daemonw/Desktop/VirtualApp.zip", "/home/daemonw/Desktop/test.zip", key[:], iv[:])
}
