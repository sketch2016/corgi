package Utils

import (
	"fmt"
	"time"
)

type logleveltype int

const (
	logLevelDebug logleveltype = 0
	logLevelInfo
	logLevelWarning
	logLevelError
)

var loglevel = logLevelDebug

//LOGD is used for log print
func LOGD(tag string, a ...interface{}) {
	if loglevel <= logLevelDebug {
		logprint(tag, a...)
	}
}

//LOGI is used for log print
func LOGI(tag string, a ...interface{}) {
	if loglevel <= logLevelDebug {
		logprint(tag, a...)
	}
}

//LOGW is used for log print
func LOGW(tag string, a ...interface{}) {
	if loglevel <= logLevelDebug {
		logprint(tag, a...)
	}
}

//LOGE is used for log print
func LOGE(tag string, a ...interface{}) {
	if loglevel <= logLevelDebug {
		logprint(tag, a...)
	}
}

func logprint(tag string, a ...interface{}) {
	t := time.Now().Format("2006-01-02 15:04:05.000")
	fmt.Print(t, " [", tag, "] ")
	fmt.Println(a...)
}
