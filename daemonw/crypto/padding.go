package crypto

//PKCS7Padding will pad the byte length, until it can
//be devisioned by the  blockSize
func PKCS7Padding(origin []byte, blockSize int) []byte {
	padSize := blockSize - len(origin)%blockSize
	padding := make([]byte, padSize)
	for i := 0; i < padSize; i++ {
		padding[i] = uint8(padSize)
	}
	return append(origin, padding...)
}

//PKCS7Padding will remove the redundant bytes of the padding
//data
func PKCS7Unpadding(origin []byte) []byte {
	originLen := len(origin)
	padSize := int(origin[originLen-1])
	return origin[:originLen-padSize]
}
