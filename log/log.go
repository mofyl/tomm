package log

import "go.uber.org/zap"

var (
	logger *zap.Logger
)

func init() {
	cfg := getDefaultLog()
	logger = buildLog(cfg)
}

func Debug(msg string, field ...zap.Field) {
	defaultLog(DEBUG, msg, field...)
}

func Msg(level LogLEVEL, msg string) {
	defaultLog(level, msg)
}

func Info(msg string, field ...zap.Field) {
	defaultLog(INFO, msg, field...)
}

func Warn(msg string, field ...zap.Field) {
	defaultLog(WARN, msg, field...)
}

func Error(msg string, field ...zap.Field) {
	defaultLog(ERROR, msg, field...)
}

func Panic(msg string, field ...zap.Field) {
	defaultLog(PANIC, msg, field...)
}

func Fatal(msg string, field ...zap.Field) {
	defaultLog(FATAL, msg, field...)
}

func defaultLog(level string, msg string, field ...zap.Field) {
	if ce := logger.Check(getLogLevel(level), msg); ce != nil {
		ce.Write(field...)
	}
}
