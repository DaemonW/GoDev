package util

import (
	"os"
	"path/filepath"
)

func ExistFile(filepath string) bool {
	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		//other error, usually permission error
		panic(err)
	}
	//no error, exist
	return true
}

func ListFilesInfo(dir string) []os.FileInfo {
	file, err := os.OpenFile(dir, os.O_RDONLY, 0444)
	if err != nil {
		return nil
	}
	defer file.Close()
	infoList, err := file.Readdir(0)
	if err != nil {
		return nil
	}
	return infoList
}

func ListFiles(dir string) []string {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil
	}
	infoList := ListFilesInfo(dir)
	if infoList == nil {
		return nil
	}
	size := len(infoList)
	subFiles := make([]string, size)
	for i := 0; i < size; i++ {
		subFiles[i] = dir + string(filepath.Separator) + infoList[i].Name()
	}
	return subFiles
}

func ListExtFiles(dir string, ext string) []string {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil
	}
	infoList := ListFilesInfo(dir)
	if infoList == nil {
		return nil
	}
	size := len(infoList)
	subFiles := make([]string, size)
	for i := 0; i < size; i++ {
		subFiles[i] = dir + string(filepath.Separator) + infoList[i].Name()
	}
	return subFiles
}
