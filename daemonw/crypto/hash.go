//util hash functions for crypto, such as md5, sha,crc32
package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"os"
	"strconv"
	"strings"
)

//GetHash caculate the byte array's hash value, return a hex base string
func GetHash(data []byte, alg string) []byte {
	var hasher hash.Hash
	alg = strings.ToUpper(alg)
	switch alg {
	case "MD5":
		hasher = md5.New()
	case "SHA256":
		hasher = sha256.New()
	case "SHA1":
		hasher = sha1.New()
	default:
		panic(errors.New("no such hash algorithm"))
	}
	hasher.Write([]byte(data))
	return hasher.Sum(nil)
}

func GetHashStr(data []byte, alg string) string {
	return Bytes2Hex(GetHash(data, alg))
}

//GetFileHash caculte the file's hash value
func GetFileHash(filepath string, alg string) ([]byte, error) {
	alg = strings.ToUpper(alg)
	var hasher hash.Hash
	switch alg {
	case "MD5":
		hasher = md5.New()
	case "SHA256":
		hasher = sha256.New()
	case "SHA1":
		hasher = sha1.New()
	default:
		err := fmt.Errorf(`algorithm '%s' is not support`, alg)
		return nil, err
	}
	file, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if _, err = io.Copy(hasher, file); err != nil {
		return nil, err
	}
	bytes := hasher.Sum(nil)
	return bytes, nil
}

func GetCrc32(data []byte) uint32 {
	c := crc32.New(crc32.IEEETable)
	c.Write(data)
	return c.Sum32()
}

func GetFileCrc32(filepath string) (uint32, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	c := crc32.New(crc32.IEEETable)
	if _, err = io.Copy(c, file); err != nil {
		return 0, err
	}
	return c.Sum32(), nil
}

//get each block's md5 with the given blocksize
func GetSplitFileMd5(filepath string, blocksize int64) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	result := make([]string, 0)
	m := md5.New()
	for {
		_, err := io.CopyN(m, file, blocksize)
		if err != nil {
			if err == io.EOF {
				result = append(result, Bytes2Hex(m.Sum(nil)))
				return result, nil
			} else {
				return nil, err
			}
		}
		result = append(result, Bytes2Hex(m.Sum(nil)))
		m.Reset()
	}
	return result, nil
}

func GetSplitFileCrc32(filepath string, blockSize int64) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	c := crc32.New(crc32.IEEETable)
	result := make([]string, 0)
	for {
		//file.Seek(total, 0)
		_, err := io.CopyN(c, file, blockSize)
		if err != nil {
			if err == io.EOF {
				result = append(result, strconv.FormatUint(uint64(c.Sum32()), 10))
				return result, nil
			} else {
				return nil, err
			}
		}
		result = append(result, strconv.FormatUint(uint64(c.Sum32()), 10))
		c.Reset()
	}
	return result, nil
}

//get md5 of the string,
func GetStrMd5(str string) string {
	if str == "" {
		return ""
	}
	bytes := []byte(str)
	m := md5.New()
	return Bytes2Hex(m.Sum(bytes))
}

//get sha256 of the string
func GetStrSha256(str string) string {
	if str == "" {
		return ""
	}
	bytes := []byte(str)
	result := sha256.Sum256(bytes)
	return Bytes2Hex(result[:])
}
