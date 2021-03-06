package xlog

import (
	"daemonw/util"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	defaultFileName   = "current.log"
	defaultLogMaxSize = 2 * 1024 * 1024
)

type fileWriter struct {
	filePath string
	time     time.Time
	f        *os.File
	size     int64
	maxSize  int64
	locker   sync.Mutex
}

func NewFileWriter(filePath string, maxSize int64) *fileWriter {
	writer := &fileWriter{filePath: filePath, time: time.Now(), maxSize: maxSize}
	writer.output(filePath)
	return writer
}

func (writer *fileWriter) Write(p []byte) (n int, err error) {
	writer.locker.Lock()
	n, err = writer.f.Write(p)
	writer.size += int64(n)
	if writer.size >= writer.maxSize {
		writer.next()
	}
	writer.locker.Unlock()
	return n, err
}

func (writer *fileWriter) next() {
	err := writer.f.Sync()
	util.PanicIfErr(err)
	err = writer.f.Close()
	util.PanicIfErr(err)
	dir := filepath.Dir(writer.filePath)
	err = os.Rename(writer.filePath, dir+string(filepath.Separator)+writer.time.Format("2006-01-02_15:04:05")+".log")
	util.PanicIfErr(err)
	writer.output(writer.filePath)
}

func (writer *fileWriter) output(filePath string) {
	var file *os.File
	fi, err := os.Stat(filePath)
	if err != nil {
		writer.size = 0
		writer.time = time.Now()
		if os.IsNotExist(err) {
			file, err = os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
		}
	} else {
		file, err = os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0666)
		writer.size = fi.Size()
		//writer.f.Seek(writer.size, 0)
	}
	if err != nil {
		panic(err)
	}
	writer.f = file
}
