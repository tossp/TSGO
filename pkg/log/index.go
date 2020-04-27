package log

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	z      *zap.SugaredLogger
	isDev  = true
	lumlog = &lumberjack.Logger{
		Filename: "",
		Compress: true,
		MaxSize:  1024 * 1024 * 100, // 单次写入容量，100MB
		MaxAge:   180,               // days
	}
)

func init() {
	Make()
	Info("日志系统启动完成")
}

//SetConfig 设置参数
func SetConfig(b bool, p string) {
	isDev = b
	os.MkdirAll(p, 0600)
	lumlog.Filename = p + "/server.log"
	Make()
}

func Make() {
	//goVersion := runtime.Version()
	//hostMachine := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	var l *zap.Logger
	if isDev {
		l, _ = zap.NewDevelopment(zap.AddCaller(), zap.AddCallerSkip(1),
			zap.AddStacktrace(zap.ErrorLevel))
	} else {
		l, _ = zap.NewProduction()
	}
	if lumlog.Filename != "" {
		l = l.WithOptions(zap.Hooks(lumberjackZapHook))
	}
	z = l.Named("site").Sugar()
	Debug("配置日志系统")
}

func lumberjackZapHook(e zapcore.Entry) error {
	lumlog.Write([]byte(fmt.Sprintf("%s    ", e.Time.Format("2006-01-02T15:04:05Z07"))))
	lumlog.Write([]byte(fmt.Sprintf("%-8s  ", e.Level.CapitalString())))
	lumlog.Write([]byte(fmt.Sprintf("%s    ", e.LoggerName)))
	lumlog.Write([]byte(fmt.Sprintf("%s    ", e.Caller.TrimmedPath())))
	lumlog.Write([]byte(fmt.Sprintf("%s    ", e.Message)))
	if e.Stack != "" {
		lumlog.Write([]byte(fmt.Sprintf("\n%s\n", e.Stack)))
	} else {
		lumlog.Write([]byte("\n"))
	}

	return nil
}
