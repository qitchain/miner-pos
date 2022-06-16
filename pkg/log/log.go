package log

import (
	"chia-miner/utils"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path"
)

const (
	logPath = "./log"
)

func InitLog(level, logFile string) {
	// default logger
	outList := make([]io.Writer, 0)
	outList = append(outList, os.Stdout)
	if logFile != "" {
		if !utils.PathExists(logPath) {
			_ = os.Mkdir(logPath, os.ModePerm)
		}
		writerFile, err := os.OpenFile(path.Join(logPath, logFile), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			log.Fatalf("Failed to create file %v error %v", logFile, err)
		} else {
			outList = append(outList, writerFile)
		}
	}
	
	logrus.SetOutput(io.MultiWriter(outList...))
	if level == "" {
		level = "info"
	}
	if level == "info" {
		logrus.SetLevel(logrus.InfoLevel)
	} else if level == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	} else if level == "error" {
		logrus.SetLevel(logrus.ErrorLevel)
	} else if level == "trace" {
		logrus.SetLevel(logrus.TraceLevel)
	} else if level == "warn" {
		logrus.SetLevel(logrus.WarnLevel)
	}
}
