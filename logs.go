package zap_wrapper

import (
	"bytes"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 配置结构体
type Config struct {
	TimeLayout string // 时间布局，默认为 "2006-01-02 15:04:05"
}

// 自定义带有缓冲区的核心
type bufferedCore struct {
	zapcore.Core
	encoder zapcore.Encoder
	buffer  *bytes.Buffer
}

// 创建一个新的bufferedCore实例
func newBufferedCore(core zapcore.Core, encoder zapcore.Encoder) *bufferedCore {
	return &bufferedCore{
		Core:    core,
		encoder: encoder,
		buffer:  new(bytes.Buffer),
	}
}

// 实现zapcore.Core接口的Write方法
func (bc *bufferedCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	encoded, err := bc.encoder.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}
	defer encoded.Free()

	_, err = fmt.Fprint(bc.buffer, encoded.String())
	return err
}

// 获取缓冲区内容并清空
func (bc *bufferedCore) GetOutputAndClear() string {
	output := bc.buffer.String()
	bc.buffer.Reset() // 清空缓冲区
	return output
}

// Logger 结构体
type Logger struct {
	cores       map[zapcore.Level]*bufferedCore
	timeEncoder zapcore.TimeEncoder
}

// 创建一个新的Logger实例
func NewLogger(config *Config) *Logger {
	if config == nil || config.TimeLayout == "" {
		config = &Config{TimeLayout: "2006-01-02 15:04:05"}
	}
	timeEncoder := zapcore.TimeEncoderOfLayout(config.TimeLayout)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = timeEncoder

	cores := make(map[zapcore.Level]*bufferedCore)
	for _, level := range []zapcore.Level{zap.DebugLevel, zap.InfoLevel, zap.WarnLevel, zap.ErrorLevel, zap.DPanicLevel, zap.PanicLevel, zap.FatalLevel} {
		encoder := zapcore.NewJSONEncoder(encoderConfig)
		core := zapcore.NewCore(encoder, zapcore.AddSync(new(bytes.Buffer)), level)
		cores[level] = newBufferedCore(core, encoder)
	}

	return &Logger{
		cores:       cores,
		timeEncoder: timeEncoder,
	}
}

// 设置时间编码器
func (l *Logger) SetTimeLayout(layout string) {
	l.timeEncoder = zapcore.TimeEncoderOfLayout(layout)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = l.timeEncoder

	for level, bc := range l.cores {
		encoder := zapcore.NewJSONEncoder(encoderConfig)
		newCore := zapcore.NewCore(encoder, zapcore.AddSync(bc.buffer), level)
		l.cores[level] = newBufferedCore(newCore, encoder)
	}
}

// 返回格式化后的日志信息作为字符串
func (l *Logger) Log(level zapcore.Level, msg string, keysAndValues ...interface{}) string {
	bc, exists := l.cores[level]
	if !exists {
		bc = l.cores[zap.InfoLevel] // 默认使用 Info 级别核心
	}

	logger := zap.New(bc).Sugar()
	switch level {
	case zap.DebugLevel:
		logger.Debugw(msg, keysAndValues...)
	case zap.InfoLevel:
		logger.Infow(msg, keysAndValues...)
	case zap.WarnLevel:
		logger.Warnw(msg, keysAndValues...)
	case zap.ErrorLevel:
		logger.Errorw(msg, keysAndValues...)
	case zap.DPanicLevel:
		logger.DPanicw(msg, keysAndValues...)
	case zap.PanicLevel:
		logger.Panicw(msg, keysAndValues...)
	case zap.FatalLevel:
		logger.Fatalw(msg, keysAndValues...)
	default:
		logger.Infow(msg, keysAndValues...) // 默认情况
	}

	// 获取缓冲区内容并清空
	return bc.GetOutputAndClear()
}

// 提供不同级别的日志记录方法
func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) string {
	return l.Log(zap.DebugLevel, msg, keysAndValues...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) string {
	return l.Log(zap.InfoLevel, msg, keysAndValues...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) string {
	return l.Log(zap.WarnLevel, msg, keysAndValues...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) string {
	return l.Log(zap.ErrorLevel, msg, keysAndValues...)
}

func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) string {
	return l.Log(zap.DPanicLevel, msg, keysAndValues...)
}

func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) string {
	return l.Log(zap.PanicLevel, msg, keysAndValues...)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) string {
	return l.Log(zap.FatalLevel, msg, keysAndValues...)
}
