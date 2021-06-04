package log

import (
	"github.com/jackc/pgconn"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/tossp/tsgo/pkg/errors"
)

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
func WarnErr(err error, args ...interface{}) {
	if err != nil {
		if len(args) > 0 {
			z.Warn(append([]interface{}{err}, args...))
		} else {
			z.Warn(err)
		}

	}
}
func WarnDB(tx *gorm.DB, args ...interface{}) (err error) {
	err = tx.Error
	if err != nil {
		zp := z.With()
		pgerr, ok := err.(*pgconn.PgError)
		if ok {
			zp = zp.With(
				"RowsAffected", tx.RowsAffected,
				"pgErr", pgerr,
			)
			z.Warnf("%#v", pgerr)
			err = errors.ErrDatabase
		}
		if len(args) > 0 {
			zp.Warn(append([]interface{}{err}, args...))
		} else {
			zp.Warn(err)
		}
	}
	return
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

func DPanic(args ...interface{}) {
	z.DPanic(args...)
}
func Panic(args ...interface{}) {
	z.Panic(args...)
}

func Fatal(args ...interface{}) {
	z.Fatal(args...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	z.Fatalw(msg, keysAndValues...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func Panicw(msg string, keysAndValues ...interface{}) {
	z.Panicw(msg, keysAndValues...)
}
