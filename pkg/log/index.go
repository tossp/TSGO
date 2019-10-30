package log

import (
	"fmt"
	"runtime"
	"strings"

	//"tsl/core/setting/machineid"

	graylog "github.com/gemnasium/logrus-graylog-hook/v3"
	"github.com/sirupsen/logrus"
)

var log = NewLog()

type logger struct {
	*logrus.Logger
}

func NewLog() (l *logger) {
	//logrus.SetSkipPackageNameForCaller("tsl/core/log")
	//logrus.SetSkipPackageNameForCaller("tsl/core/db/logger")
	//logrus.SetSkipPackageNameForCaller("github.com/go-xorm/xorm")
	//logrus.SetSkipPackageNameForCaller("github.com/jinzhu/gorm")
	goVersion := runtime.Version()
	hostMachine := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	//TODO 需要修改输出格式
	// fmt.Fprintf(b, "\x1b[%dm%s[%04d]%s\x1b[0m\n%-44s ", levelColor, levelText, int(entry.Time.Sub(baseTimestamp)/time.Second), caller, entry.Message)
	formatter := &logrus.TextFormatter{
		//FullTimestamp: true,
		ForceColors: true,
		FieldMap: logrus.FieldMap{
			"FieldKeyTime":  "@timestamp",
			"FieldKeyLevel": "@level",
			"FieldKeyMsg":   "@message",
		},
	}
	hook := graylog.NewGraylogHook("log.tossp.com:12201", map[string]interface{}{
		"app":         "tsl",
		"golang":      goVersion,
		"hostMachine": hostMachine,
		//"MachineID":   machineid.Id(),
	})
	l = &logger{
		Logger: logrus.New(),
	}
	l.SetFormatter(formatter)
	//l.SetReportCaller(true)
	l.SetLevel(logrus.TraceLevel)
	l.AddHook(hook)
	l.Info("日志系统启动完成")
	return
}

func GetLogger() *logger {
	return log
}

func SetReportCaller(reportCaller bool) {
	log.SetReportCaller(reportCaller)
}
func SetLevel(lvl string) {
	lvl = strings.ToLower(lvl)
	if lvl == "trace" {
		lvl = "debug"
	}
	if lvl == "code@tossp.com" {
		lvl = "trace"
	}
	lv, err := logrus.ParseLevel(lvl)
	if err != nil {
		log.Warn("日志等级关键字错误，加载默认等级")
		lv = logrus.WarnLevel
	}
	SetReportCaller(lv == logrus.TraceLevel)
	log.SetLevel(lv)
}
