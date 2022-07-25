package logger

import (
	"go.uber.org/zap"
)

func LogInit() {
	l, _ := zap.NewProduction(zap.AddCallerSkip(1))
	defer l.Sync()

	zap.ReplaceGlobals(l)
}

func With(fields ...zap.Field) *zap.Logger {
	return zap.L().WithOptions(zap.AddCallerSkip(-1)).With(fields...)
}

func Debug(msg string, fields ...zap.Field) {
	zap.L().Debug(msg, fields...)
}

func Debugf(template string, args ...interface{}) {
	zap.S().Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	zap.S().Debugw(msg, keysAndValues...)
}

func Info(msg string, fields ...zap.Field) {
	zap.L().Info(msg, fields...)
}

func Infof(template string, args ...interface{}) {
	zap.S().Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	zap.S().Infow(msg, keysAndValues...)
}

func Warn(msg string, fields ...zap.Field) {
	zap.L().Warn(msg, fields...)
}

func Warnf(template string, args ...interface{}) {
	zap.S().Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	zap.S().Warnw(msg, keysAndValues...)
}

func Error(msg string, fields ...zap.Field) {
	zap.L().Error(msg, fields...)
}

func Errorf(template string, args ...interface{}) {
	zap.S().Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	zap.S().Errorw(msg, keysAndValues...)
}
