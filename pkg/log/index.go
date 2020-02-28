package log

import (
	"go.uber.org/zap"
)

var (
	z     *zap.SugaredLogger
	isDev = true
)

func init() {
	Make()
	Info("日志系统启动完成")
}

func SetMode(b bool) {
	isDev = b
	Make()
}

func Make() {
	//goVersion := runtime.Version()
	//hostMachine := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	var l *zap.Logger
	if isDev {
		l, _ = zap.NewDevelopment(zap.AddCaller(), zap.AddCallerSkip(1))
	} else {
		l, _ = zap.NewProduction(zap.AddCaller(), zap.AddCallerSkip(1))
	}
	z = l.Named("site").Sugar()
	Debug("配置日志系统")
}
