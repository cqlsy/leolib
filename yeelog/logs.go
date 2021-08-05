/**
 * Created by angelina-zf on 17/2/25.
 */

// yeego 日志相关的功能
// 依赖： "github.com/Sirupsen/logrus"
// 基于logrus的log类,自己搞成了分级的形式
// 可以设置将error层次的log发送到es
package yeelog

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"lib/yeefile"
	"lib/yeetime"
	"os"
	"runtime"
	"sync"
	"time"
)

var day *SafeMapDay
var logfile *SafeMapFile
var loggers *SafeMapLogger

// logPath
// 日志的存储地址，按照时间存储
var logPath string

// timePath
// 日志的存储地址，按照时间存储
var timePath string = yeetime.DateFormat(time.Now(), "YYYY-MM-DD") + "/"

// runMode
// 运行环境 默认为dev
var runMode string = "dev"

// MustInitLogs
// 注册log
// @param logpath 日志位置
// @param runmode 运行环境 dev|pro
func MustInitLog(path, mode string) {
	if mode != "" && (mode == "dev" || mode == "pro") {
		runMode = mode
	}
	logPath = path
	if runMode != "dev" {
		if !yeefile.FileExists(logPath) {
			if createErr := os.MkdirAll(logPath, os.ModePerm); createErr != nil {
				panic("error to create logs path : " + createErr.Error())
			}
		}
	}
	day = newSafeMapDay()
	logfile = newSafeMapFile()
	loggers = newSafeMapLogger()
}

func initLogger(logFile string) {
	loggers.writeMap(logFile, logrus.New())
	setLogSConfig(logFile)
}

func setLogSConfig(level string) {
	//var err error
	logger := loggers.readMap(level)
	logger.Formatter = new(logrus.JSONFormatter)
	if runMode != "dev" {
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
	}
	if runMode == "dev" {
		// 加上这个貌似没颜色了~好奇怪啊！！！
		logger.Out = os.Stdout
	} else {
		logger.Out = logfile.readMap(level)
	}
	day.writeMap(level, yeetime.Day())

}

// 检查日志，如果不是当天的日志文件则进行更新
func checkAndUpdateLogFile(logFile string) {
	day2 := yeetime.Day()
	if day2 != day.readMap(logFile) {
		defer logfile.readMap(logFile).Close()
		timePath = yeetime.DateFormat(time.Now(), "YYYY-MM-DD") + "/" // 更新time对应的文件夹
		if yeefile.FileExists(getLogFullPath(logFile)) {
			file, err := os.OpenFile(getLogFullPath(logFile), os.O_RDWR|os.O_APPEND, 0660)
			if err == nil {
				logfile.writeMap(logFile, file)
			} else {
				Print(err)
			}
		} else {
			file, err := os.Create(getLogFullPath(logFile))
			if err == nil {
				logfile.writeMap(logFile, file)
			} else {
				Print(err)
			}
		}
		day.writeMap(logFile, day2)
		loggers.readMap(logFile).Out = logfile.readMap(logFile)
	}
}

// locate
// 找到是哪个文件的哪个地方打出的log
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

// LogInfo
// 记录Info信息
func LogDefault(str interface{}, logFile string) {
	Log(str, LogFields{}, logFile)
}

// LogInfo
// 记录Info信息
func Log(str interface{}, data LogFields, logFile string) {
	if loggers.readMap(logFile) == nil {
		initLogger(logFile)
	}
	if runMode != "dev" { // 更新文件
		checkAndUpdateLogFile(logFile)
	}
	loggers.readMap(logFile).WithFields(logrus.Fields(locate(data))).Info(str)
}

func LogDebugDefault(str interface{}) {
	LogDebug(str, LogFields{})
}

func LogDebug(str interface{}, data LogFields) {
	logFile := logrus.DebugLevel.String()
	if loggers.readMap(logFile) == nil {
		initLogger(logFile)
	}
	if runMode != "dev" { // 更新文件
		checkAndUpdateLogFile(logFile)
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
	if runMode != "dev" { // 更新文件
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
	if runMode != "dev" { // 更新文件
		checkAndUpdateLogFile(logFile)
	}
	loggers.readMap(logFile).WithFields(logrus.Fields(locate(data))).Error(str)
}

// 检查文件是否过期了 circle 天数
func CheckExpiredLog(circle int64) {
	if yeefile.FileExists(logPath) { // 检查是否有当前的目录
		files, _ := ioutil.ReadDir(logPath)
		now := time.Now().Unix()
		for _, f := range files {
			t, err := time.Parse("2006-01-02", f.Name())
			if err == nil {
				if (now - t.Unix()) > 60*60*24*circle {
					err := os.RemoveAll(logPath + "/" + f.Name())
					if err != nil {
						LogDefault("Delete ExpiredLog Err: "+err.Error(), "log")
					}
				}
			}
		}
	}
}

type LogFields logrus.Fields

type SafeMapDay struct {
	sync.RWMutex
	Map map[string]int
}

func newSafeMapDay() *SafeMapDay {
	sm := new(SafeMapDay)
	sm.Map = make(map[string]int)
	return sm
}

func (sm *SafeMapDay) readMap(key string) int {
	sm.RLock()
	value := sm.Map[key]
	sm.RUnlock()
	return value
}

func (sm *SafeMapDay) writeMap(key string, value int) {
	sm.Lock()
	sm.Map[key] = value
	sm.Unlock()
}

type SafeMapFile struct {
	sync.RWMutex
	Map map[string]*os.File
}

func newSafeMapFile() *SafeMapFile {
	sm := new(SafeMapFile)
	sm.Map = make(map[string]*os.File)
	return sm
}

func (sm *SafeMapFile) readMap(key string) *os.File {
	sm.RLock()
	value := sm.Map[key]
	sm.RUnlock()
	return value
}

func (sm *SafeMapFile) writeMap(key string, value *os.File) {
	sm.Lock()
	sm.Map[key] = value
	sm.Unlock()
}

type SafeMapLogger struct {
	sync.RWMutex
	Map map[string]*logrus.Logger
}

func newSafeMapLogger() *SafeMapLogger {
	sm := new(SafeMapLogger)
	sm.Map = make(map[string]*logrus.Logger)
	return sm
}

func (sm *SafeMapLogger) readMap(key string) *logrus.Logger {
	sm.RLock()
	value := sm.Map[key]
	sm.RUnlock()
	return value
}

func (sm *SafeMapLogger) writeMap(key string, value *logrus.Logger) {
	sm.Lock()
	sm.Map[key] = value
	sm.Unlock()
}
