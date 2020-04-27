package log

import "go.uber.org/zap"

func Desugar() *zap.Logger {
	return z.Desugar()
}

func Logger() *zap.SugaredLogger {
	return z
}
func With(args ...interface{}) *zap.SugaredLogger {
	return z.With(args...)
}
func Debugf(format string, args ...interface{}) {
	z.Debugf(format, args...)
}

// Infof formats according to a format specifier and returns the resulting string.
func Infof(format string, args ...interface{}) {
	z.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	z.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	z.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	z.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	z.Panicf(format, args...)
}

func Debug(args ...interface{}) {
	z.Debug(args...)
}
func Debugw(msg string, keysAndValues ...interface{}) {
	z.Debugw(msg, keysAndValues...)
}

func Info(args ...interface{}) {
	z.Info(args...)
}
func Infow(msg string, keysAndValues ...interface{}) {
	z.Infow(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	z.Warn(args...)
}
func Warnw(msg string, keysAndValues ...interface{}) {
	z.Warnw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	z.Error(args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	z.Errorw(msg, keysAndValues...)
}

func Fatal(args ...interface{}) {
	z.Fatal(args...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	z.Fatalw(msg, keysAndValues...)
}

func Panic(args ...interface{}) {
	z.Panic(args...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func Panicw(msg string, keysAndValues ...interface{}) {
	z.Panicw(msg, keysAndValues...)
}
