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
	logLevel = ""
)

func init() {

}

func make() {
	var syncer zapcore.WriteSyncer
	syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumlog))
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "line",
		MessageKey:       "msg",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:       zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration:   zapcore.SecondsDurationEncoder, //
		EncodeCaller:     zapcore.ShortCallerEncoder,     // 全路径编码器
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: "\t",
	}
	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	var encoder zapcore.Encoder
	encoder = zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(
		encoder,
		syncer,
		level,
	)
	//goVersion := runtime.Version()
	//hostMachine := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	var l *zap.Logger
	l = zap.New(core)
	if isDev {
		l = l.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
		// l, _ = zap.NewDevelopment(zap.AddCaller(), zap.AddCallerSkip(1),
		// 	zap.AddStacktrace(zap.ErrorLevel))
	}
	l.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, NewSentryCore(SentryCoreConfig{
			Level: zap.WarnLevel,
			//Tags: map[string]string{
			//    "source": "demo",
			//},
		}))
	}))
	// if lumlog.Filename != "" {
	// 	l = l.WithOptions(zap.Hooks(lumberjackZapHook))
	// }
	z = l.Named("admission").Sugar()
	//Debug("日志系统重载完成")
}

//SetConfig 设置参数
func SetConfig(b bool, p, lv string) {
	isDev = b
	os.MkdirAll(p, 0600)
	lumlog.Filename = p + "/server.log"
	logLevel = lv
	make()
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
