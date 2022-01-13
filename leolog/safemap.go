package leolog

import (
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

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

