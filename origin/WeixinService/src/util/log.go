package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

// GetLogger getLogger
func GetLogger(fileName string, prefix string) *log.Logger {
	// 定义一个文件
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		os.MkdirAll(path.Dir(fileName), 0777)
		logFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln(fmt.Sprintf("open log file: %s faild", fileName))
			return nil
		}
	}
	// 创建一个日志对象，同时打到系统和文件中
	logger := log.New(io.MultiWriter(os.Stderr, logFile), prefix, log.LstdFlags)
	//配置log的Flag参数
	logger.SetFlags(logger.Flags() | log.LstdFlags | log.Llongfile)
	return logger
}
