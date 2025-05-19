package zap_wrapper

import (
	"go.uber.org/zap"
)

type Logger struct {
	sugar *zap.SugaredLogger
}

// 初始化函数，可传入 zap.Logger 配置
func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{
		sugar: logger.Sugar(),
	}
}

// 获取底层 SugaredLogger（可选）
func (l *Logger) Desugar() *zap.SugaredLogger {
	return l.sugar
}

func NewLoggerProduction() *Logger {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	return &Logger{
		sugar: logger.Sugar(),
	}
}

func NewLoggerDevelopment() *Logger {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	return &Logger{
		sugar: logger.Sugar(),
	}
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) string {
	return l.sugar.Debugw2(msg, keysAndValues...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) string {
	return l.sugar.Infow2(msg, keysAndValues...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) string {
	return l.sugar.Warnw2(msg, keysAndValues...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) string {
	return l.sugar.Errorw2(msg, keysAndValues...)
}

func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) string {
	return l.sugar.Panicw2(msg, keysAndValues...)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) string {
	return l.sugar.Fatalw2(msg, keysAndValues...)
}
