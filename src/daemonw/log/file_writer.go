package log

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
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
}

func NewFileWriter(filePath string, maxSize int64) *fileWriter {
	writer := &fileWriter{filePath: filePath, time: time.Now(), maxSize: maxSize}
	fi, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			writer.f, err = os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
			writer.Write([]byte(writer.time.Format("2006-01-02 15:04:05") + "\n"))
		}
	} else {
		writer.f, err = os.OpenFile(filePath, os.O_RDWR, 0666)
		writer.size = fi.Size()
		var timeStamp []byte
		timeStamp, _, err = bufio.NewReader(writer.f).ReadLine()
		if timeStamp != nil && len(timeStamp) > 0 {
			writer.time, err = time.Parse("2006-01-02 15:04:05", string(timeStamp))
		}
		writer.f.Seek(writer.size, 0)
	}
	if err != nil {
		log.Fatal(err)
	}
	return writer
}

func (writer *fileWriter) Write(p []byte) (n int, err error) {
	n, err = writer.f.Write(p)
	writer.size += int64(n)
	if writer.size >= writer.maxSize {
		writer.next()
	}
	return
}

func (writer *fileWriter) next() {
	writer.f.Sync()
	writer.f.Close()
	dir := filepath.Dir(writer.filePath)
	os.Rename(writer.filePath, dir+string(filepath.Separator)+writer.time.Format("2006-01-02_15:04:05")+".log")
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
			file.Write([]byte(writer.time.Format("2006-01-02 15:04:05") + "\n"))
		}
	} else {
		file, err = os.OpenFile(filePath, os.O_RDWR, 0666)
		writer.size = fi.Size()
		var timeStamp []byte
		timeStamp, _, err = bufio.NewReader(writer.f).ReadLine()
		if timeStamp != nil && len(timeStamp) > 0 {
			writer.time, err = time.Parse("2006-01-02 15:04:05", string(timeStamp))
		}
		writer.f.Seek(writer.size, 0)
	}
	if err != nil {
		log.Fatal(err)
	}
	writer.f = file
}
