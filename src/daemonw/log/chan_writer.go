package log

import (
	"time"
	"os"
	"path/filepath"
	"sync"
	"fmt"
	dlog "log"
)

type ChanWriter struct {
	filePath string
	msgChan  chan []byte
	cap      int
	time     time.Time
	f        *os.File
	size     int64
	maxSize  int64
	locker   sync.Locker
}

func NewChanWriter(filePath string, maxSize int64) *ChanWriter {
	writer := &ChanWriter{filePath: filePath, time: time.Now(), maxSize: maxSize}
	writer.cap = 1024 * 16
	writer.msgChan = make(chan []byte, writer.cap)
	writer.output(filePath)
	return writer
}

func (writer *ChanWriter) Write(p []byte) (n int, err error) {
	writer.msgChan <- p
	return len(p), nil
}

func (writer *ChanWriter) Sync() {
	go func() {
		for {
			msg, ok := <-writer.msgChan
			if !ok {
				break
			}
			n, err := writer.f.Write(msg)
			if err != nil {
				fmt.Println(err)
			}
			writer.size += int64(n)
			if writer.size >= writer.maxSize {
				writer.next()
			}
		}
	}()
}

func (writer *ChanWriter) next() {
	writer.f.Sync()
	writer.f.Close()
	dir := filepath.Dir(writer.filePath)
	os.Rename(writer.filePath, dir+string(filepath.Separator)+writer.time.Format("2006-01-02_15:04:05")+".log")
	writer.output(writer.filePath)
}

func (writer *ChanWriter) output(filePath string) {
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
		dlog.Fatal(err)
	}
	writer.f = file
}
