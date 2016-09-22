//import "github.com/solomonooo/mercury"
//author : solomonooo
//create time : 2016-09-08

package mercury

import (
	"flag"
	"github.com/solomonooo/glog-go"
	"os"
	"strconv"
	"strings"
)

const (
	FATAL = iota
	ERROR
	WARN
	INFO
	DEBUG
)

const (
	DEFAULT_LOG_DIR   = "./log/"
	DEFAULT_LOG_LEVEL = "info"
)

func init() {
	setLogDir(DEFAULT_LOG_DIR)
	setLogLevel(DEFAULT_LOG_LEVEL)
}

func setLogDir(dir string) {
	os.MkdirAll(dir, os.ModePerm)
	flag.Set("log_dir", dir)
}

func setLogLevel(level string) {
	var v string
	switch strings.ToLower(level) {
	case "debug":
		v = strconv.Itoa(DEBUG)
	case "info":
		v = strconv.Itoa(INFO)
	case "warn":
		v = strconv.Itoa(WARN)
	case "error":
		v = strconv.Itoa(ERROR)
	case "fatal":
		v = strconv.Itoa(FATAL)
	default:
		v = strconv.Itoa(INFO)
	}

	flag.Set("v", v)
}

func Debug(format string, args ...interface{}) {
	defer glog.Flush()
	glog.V(DEBUG).InfofDepth(1, format, args...)
}

func Info(format string, args ...interface{}) {
	defer glog.Flush()
	glog.V(INFO).InfofDepth(1, format, args...)
}

func Warn(format string, args ...interface{}) {
	defer glog.Flush()
	glog.V(WARN).InfofDepth(1, format, args...)
}

func Error(format string, args ...interface{}) {
	defer glog.Flush()
	glog.V(ERROR).InfofDepth(1, format, args...)
}

func Fatal(format string, args ...interface{}) {
	defer glog.Flush()
	glog.V(FATAL).InfofDepth(1, format, args...)
}
