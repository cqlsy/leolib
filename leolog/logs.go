package leolog

import (
	"bytes"
	"fmt"
	"github.com/leo/leolib/leofile"
	"github.com/leo/leolib/leotime"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	PRO = "PRO"
	DEV = "DEV"
)

var (
	day      *SafeMapDay
	logfile  *SafeMapFile
	loggers  *SafeMapLogger
	logPath  string
	timePath string = leotime.DateFormat(time.Now(), "YYYY-MM-DD") + "/"
	runMode  string = DEV
)

type LogFields logrus.Fields

// MustInitLogs
// mode: dev|pro
func MustInitLog(path, mode string) {
	logPath = path
	if strings.ToUpper(mode) != PRO { // 只要不是pro，则是dev环境
		runMode = DEV
	} else {
		runMode = PRO
		if !leofile.FileExists(logPath) {
			if createErr := os.MkdirAll(logPath, os.ModePerm); createErr != nil {
				panic("error to create logs path : " + createErr.Error())
			}
		}
	}
	day = newSafeMapDay()
	logfile = newSafeMapFile()
	loggers = newSafeMapLogger()
}

func initLogger(level string) {
	loggers.writeMap(level, logrus.New())
	//var err error
	logger := loggers.readMap(level)
	logger.Formatter = new(logrus.JSONFormatter)
	if runMode == PRO {
		file, err := os.OpenFile(getLogFullPath(level), os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			file, err = os.Create(getLogFullPath(level))
			if err != nil {
				panic("error to create logs path : " + err.Error())
			} else {
				logfile.writeMap(level, file)
			}
		} else {
			logfile.writeMap(level, file)
		}
		logger.Out = logfile.readMap(level)
	} else {
		logger.Out = os.Stdout
	}
	day.writeMap(level, leotime.Day())
}

// LogInfo
// record int logFile
// isRecord
func LogCustomFile(str interface{}, logFile string, isRecord bool) {
	if loggers.readMap(logFile) == nil {
		initLogger(logFile)
	}
	if isRecord {
		checkAndUpdateLogFile(logFile)
	}
	loggers.readMap(logFile).WithFields(logrus.Fields(locate(LogFields{}))).Info(str)
}

func LogDebugDefault(str interface{}) {
	LogDebug(str, LogFields{})
}

func LogDebug(str interface{}, data LogFields) {
	if runMode != DEV {
		return
	}
	logFile := logrus.DebugLevel.String()
	if loggers.readMap(logFile) == nil {
		initLogger(logFile)
	}
	loggers.readMap(logFile).WithFields(logrus.Fields(locate(data))).Debug(str)
}

func LogInfoDefault(str interface{}) {
	LogInfo(str, LogFields{})
}

func LogInfo(str interface{}, data LogFields) {
	logFile := logrus.InfoLevel.String()
	if loggers.readMap(logFile) == nil {
		initLogger(logFile)
	}
	if runMode != DEV {
		checkAndUpdateLogFile(logFile)
	}
	loggers.readMap(logFile).WithFields(logrus.Fields(locate(data))).Info(str)
}

func LogErrorDefault(str interface{}) {
	LogError(str, LogFields{})
}

func LogError(str interface{}, data LogFields) {
	logFile := logrus.ErrorLevel.String()
	if loggers.readMap(logFile) == nil {
		initLogger(logFile)
	}
	if runMode != DEV {
		checkAndUpdateLogFile(logFile)
	}
	loggers.readMap(logFile).WithFields(logrus.Fields(locate(data))).Error(str)
}

// Print
// 格式化打印数据.
func Print(objList ...interface{}) {
	defer func() {
		_ = recover()
	}()
	var pc, file, line, ok = runtime.Caller(1)
	if !ok {
		return
	}
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			file = file[i+1:]
			break
		}
	}
	var buf = new(bytes.Buffer)
	for i := 0; i < len(objList); i++ {
		fmt.Fprintf(buf, "%c[0;0;32m [%s][%s:%d] at %s(): ", 0x1B,
			leotime.DateFormat(time.Now(), leotime.ForMate_yyyymmddhhmmss), file, line, function(pc))
		formatData(buf, objList[i])
	}
	fmt.Fprintf(buf, "%c[0m", 0x1B)
	os.Stdout.WriteString(buf.String())
}

// 只是再DEBUG环境下打印
func PrintDebug(objList ...interface{}) {
	if runMode == PRO {
		return
	}
	Print(objList)
}

// delete logs files
func DeleteExpiredLog(circle int64) {
	if leofile.FileExists(logPath) {
		files, _ := ioutil.ReadDir(logPath)
		now := time.Now().Unix()
		for _, f := range files {
			t, err := time.Parse("2006-01-02", f.Name())
			if err == nil {
				if (now - t.Unix()) > 60*60*24*circle {
					err := os.RemoveAll(logPath + "/" + f.Name())
					if err != nil {
						LogErrorDefault("Delete ExpiredLog Err: " + err.Error())
					}
				}
			}
		}
	}
}

func checkAndUpdateLogFile(logFile string) {
	day2 := leotime.Day()
	if day2 != day.readMap(logFile) {
		defer logfile.readMap(logFile).Close()
		timePath = leotime.DateFormat(time.Now(), "YYYY-MM-DD") + "/"
		if leofile.FileExists(getLogFullPath(logFile)) {
			file, err := os.OpenFile(getLogFullPath(logFile), os.O_RDWR|os.O_APPEND, 0660)
			if err == nil {
				logfile.writeMap(logFile, file)
			} else {
				Print("Log", err)
			}
		} else {
			file, err := os.Create(getLogFullPath(logFile))
			if err == nil {
				logfile.writeMap(logFile, file)
			} else {
				Print("Log", err)
			}
		}
		day.writeMap(logFile, day2)
		loggers.readMap(logFile).Out = logfile.readMap(logFile)
	}
}

// locate
func locate(fields LogFields) LogFields {
	_, path, line, ok := runtime.Caller(3)
	if ok {
		fields["file"] = path
		fields["line"] = line
	}
	return fields
}

func getLogFullPath(logFile string) string {
	os.MkdirAll(logPath+"/"+timePath, os.ModePerm)
	return logPath + "/" + timePath + logFile + ".log"
}
