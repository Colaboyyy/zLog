package zlogger

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"
)

type ConsoleZLogger struct {
	Level zLogLevel
}

// NewConsoleZLogger 构造方法
func NewConsoleZLogger(ls string) ConsoleZLogger {
	level, err := parseLogLevel(ls)
	if err != nil {
		panic(err)
	}
	//把解析得到的level赋值给结构体的Level变量
	return ConsoleZLogger{
		Level: level,
	}
}

// IsEnable 判断是否需要记录该日志
func (c *ConsoleZLogger) isEnable(level zLogLevel) bool {
	return level > c.Level
}

// getInfo 获取函数名,文件名,行号
func getInfo(skip int) (funcName, fileName string, lineNum int) {
	pc, file, lineNum, ok := runtime.Caller(skip)
	if !ok {
		fmt.Println("runtime.Caller() failed")
	}
	funcName = runtime.FuncForPC(pc).Name()
	fileName = path.Base(file)
	funcName = strings.Split(funcName, ".")[1]
	return
}

//具体记录日志的方法
func (c *ConsoleZLogger) zlog(lv zLogLevel, msg string, arg ...interface{}) {
	fullMsg := fmt.Sprintf(msg, arg...) //Sprintf根据参数生成格式化的字符串并返回该字符串
	fmt.Println(fullMsg)
	now := time.Now()
	// main > DEBUG(假設) > log 所以有3層
	funcName, fileName, lineNum := getInfo(3)
	fmt.Printf("[%s] [%s] [%s:%s:%d] %s\n", now.Format("2006-01-02 15:04:05"), getLogString(lv), White(funcName), White(fileName), lineNum, Grey(fullMsg))
}

// Debug 方法
func (c *ConsoleZLogger) Debug(msg string, arg ...interface{}) {
	if c.isEnable(DEBUG) {
		c.zlog(DEBUG, msg, arg...)
	}
}

// Info Debug 方法
func (c *ConsoleZLogger) Info(msg string, arg ...interface{}) {
	if c.isEnable(INFO) {
		c.zlog(INFO, msg, arg...)
	}
}

// Warn 方法
func (c *ConsoleZLogger) Warn(msg string, arg ...interface{}) {
	if c.isEnable(WARN) {
		c.zlog(WARN, msg, arg...)
	}
}

// Error 方法
func (c *ConsoleZLogger) Error(msg string, arg ...interface{}) {
	if c.isEnable(ERROR) {
		c.zlog(ERROR, msg, arg...)
	}
}

// Fatal 方法
func (c *ConsoleZLogger) Fatal(msg string, arg ...interface{}) {
	if c.isEnable(FATAL) {
		c.zlog(FATAL, msg, arg...)
	}
}
