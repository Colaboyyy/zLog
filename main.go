package main

import (
	"time"
	"zlog/zlogger"
)

var (
	fileName    string = "file_logger.log"
	filePath    string = "/Users/tommyyzhang/GolandProjects/zlog/"
	maxFileSize int    = 10 * 1024 * 1024
)

func main() {

	//zlog := zlogger.NewConsoleZLogger("Info")
	zlog := zlogger.NewFileZLogger("Info", fileName, filePath, maxFileSize)
	for {
		zlog.Debug("debug log...")
		zlog.Info("info log...")
		zlog.Warn("warn log...")
		id := 10010
		zlog.Error("error log: id=%d", id)
		zlog.Fatal("fatal log...")
		time.Sleep(3 * time.Second)
	}
}
