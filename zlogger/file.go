package zlogger

import (
	"fmt"
	"os"
	"path"
	"time"
)

var maxChanSize int = 50000

type MsgType struct {
	Time     string
	zLevel   zLogLevel
	FileName string
	FuncName string
	LineNum  int
	Msg      string
}

// 日志对象
type zFileLogger struct {
	// 用于限制输出层级，大于该层级才输出日志
	zLevel   zLogLevel
	FileName string
	FilePath string
	FileObj  *os.File
	// 错误文件obj
	ErrorFileObj *os.File
	// 文件最大大小
	MaxFileSize int
	// 信息管道
	MsgChan chan *MsgType
}

func NewFileZLogger(zLevel, fileName, filePath string, maxFileSize int) *zFileLogger {
	level, err := parseLogLevel(zLevel)
	if err != nil {
		panic(err)
	}
	logger := &zFileLogger{
		zLevel:      level,
		FileName:    fileName,
		FilePath:    filePath,
		MaxFileSize: maxFileSize,
		MsgChan:     make(chan *MsgType, maxChanSize),
	}
	// 初始日志文件
	err = logger.initFile()
	if err != nil {
		panic(err)
	}
	return logger
}

func (f *zFileLogger) isEnable(level zLogLevel) bool {
	return f.zLevel <= level
}

func (f *zFileLogger) initFile() error {
	file, err := os.OpenFile(path.Join(f.FilePath, f.FileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("file| os.OpenFile failed, err: %v", err)
	}
	errFile, err := os.OpenFile(path.Join(f.FilePath, f.FileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("errFile| os.OpenFile failed, err: %v", err)
	}
	f.FileObj = file
	f.ErrorFileObj = errFile
	// 后台启动goroutine等待写文件
	go f.writeFileBackend()
	return nil
}

// 检查日志大小
func (f *zFileLogger) checkFileSize(isErrorFile bool) (bool, error) {
	// 看文件是哪一种
	file := f.FileObj
	if isErrorFile {
		file = f.ErrorFileObj
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return false, fmt.Errorf("file.Stat() error: %v", err)
	}
	return int(fileInfo.Size()) <= f.MaxFileSize, nil
}

// divideFile 分割文件
func (f *zFileLogger) divideFile(isErrorFile bool) error {
	file := f.FileObj
	if isErrorFile {
		file = f.ErrorFileObj
	}
	// 把旧文件改名
	backPath := path.Join(f.FilePath, file.Name()+time.Now().Format("20060102150405000"))
	// 新增新文件与旧文件同名
	originPath := path.Join(f.FilePath, file.Name())
	// 关闭当前文件
	file.Close()
	// 将旧文件改名
	if err := os.Rename(originPath, backPath); err != nil {
		return fmt.Errorf("os.Rename() error : %v", err)
	}

	// 新开文件
	if newFile, err := os.OpenFile(originPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644); err != nil {
		return fmt.Errorf("os.OpenFile error : %v", err)
	} else {
		// 修改实例
		if isErrorFile {
			f.ErrorFileObj = newFile
		} else {
			f.FileObj = newFile
		}
	}
	return nil
}

func (f *zFileLogger) writeFileBackend() {
	for {
		select {
		case msg := <-f.MsgChan:
			// 检查日志大小
			if isValid, err := f.checkFileSize(false); err != nil {
				panic(err)
			} else if !isValid {
				f.divideFile(false)
			}
			levelStr, err := parseLogLevelToStr(msg.zLevel)
			if err != nil {
				panic(err)
			}
			//写入一般日志
			fmt.Fprintf(f.ErrorFileObj, "[%v] [%v] [%v:%v:%v] %v\n", msg.Time, levelStr, msg.FileName, msg.FuncName, msg.LineNum, msg.Msg)

			// error级别以上的日志，记录在errFile
			if msg.zLevel >= ERROR {
				// 检查日志大小
				if isValid, err := f.checkFileSize(true); err != nil {
					panic(err)
				} else if !isValid {
					f.divideFile(true)
				}
				fmt.Fprintf(f.ErrorFileObj, "[%v] [%v] [%v:%v:%v] %v\n", msg.Time, levelStr, msg.FileName, msg.FuncName, msg.LineNum, msg.Msg)
			}
		default:
			// 没有信息，休息一下
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (f *zFileLogger) zlog(level zLogLevel, msg string, arg ...interface{}) {
	if f.isEnable(level) {
		// main > DEBUG > log
		fileName, funcName, lineNum := getInfo(3)
		// 传入参数自定义
		fullMsg := fmt.Sprintf(msg, arg...)
		fmt.Println(fullMsg)

		// 将信息传入管道
		msgInstance := &MsgType{
			Time:     time.Now().Format("2006-01-02 15:04:05"),
			zLevel:   level,
			FileName: fileName,
			FuncName: funcName,
			LineNum:  lineNum,
			Msg:      fullMsg,
		}
		select {
		case f.MsgChan <- msgInstance:
		default:
			// 数据丢失的情况
			fmt.Printf("msg loss: %v\n", *msgInstance)
		}
	}
}

// Debug 方法
func (f *zFileLogger) Debug(msg string, arg ...interface{}) {
	if f.isEnable(DEBUG) {
		f.zlog(DEBUG, msg, arg...)
	}
}

// Info Debug 方法
func (f *zFileLogger) Info(msg string, arg ...interface{}) {
	if f.isEnable(INFO) {
		f.zlog(INFO, msg, arg...)
	}
}

// Warn 方法
func (f *zFileLogger) Warn(msg string, arg ...interface{}) {
	if f.isEnable(WARN) {
		f.zlog(WARN, msg, arg...)
	}
}

// Error 方法
func (f *zFileLogger) Error(msg string, arg ...interface{}) {
	if f.isEnable(ERROR) {
		f.zlog(ERROR, msg, arg...)
	}
}

// Fatal 方法
func (f *zFileLogger) Fatal(msg string, arg ...interface{}) {
	if f.isEnable(FATAL) {
		f.zlog(FATAL, msg, arg...)
	}
}
