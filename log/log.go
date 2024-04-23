package log

import (
	"io"
	"os"

	// "github.com/mattn/go-colorable"
	"gopkg.in/yaml.v3"

	// log "github.com/sirupsen/logrus"
	"github.com/logrusorgru/aurora/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	consoleConfig = zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("15:04:05.000Z0700"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     NameEncoder,
		NewReflectedEncoder: func(w io.Writer) zapcore.ReflectedEncoder {
			return yaml.NewEncoder(w)
		},
	}
	fileConfig = zapcore.EncoderConfig{
		// ... 其他配置保持不变 ...
		TimeKey:     "T",
		LevelKey:    "L",
		MessageKey:  "M",
		EncodeLevel: zapcore.CapitalLevelEncoder, // 文件使用不带颜色的编码器
		EncodeTime:  zapcore.TimeEncoderOfLayout("15:04:05.000Z0700"),
	}
	consoleEncoder = zapcore.NewConsoleEncoder(consoleConfig)
	fileEncoder    = zapcore.NewConsoleEncoder(fileConfig)
)
var LogLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
var Trace bool

var logger = zap.New(zapcore.NewTee(
	zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), LogLevel),
	zapcore.NewCore(fileEncoder, zapcore.AddSync(multipleWriter), LogLevel),
))
var sugaredLogger = logger.Sugar()
var LocaleLogger *Logger

func NameEncoder(loggerName string, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(aurora.Colorize(loggerName, aurora.WhiteFg|aurora.BlackBg).String())
}

type Zap interface {
	Lang(lang map[string]string) *Logger
	Named(name string) *Logger
	With(fields ...zap.Field) *Logger
	Trace(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

type Logger struct {
	*zap.Logger
	lang map[string]string
}

func (l Logger) Lang(lang map[string]string) *Logger {
	l.Logger = logger
	l.lang = lang
	return &l
}

func (l Logger) Named(name string) *Logger {
	l.Logger = l.Logger.Named(name)
	return &l
}

func (l Logger) With(fields ...zap.Field) *Logger {
	for i, field := range fields {
		if v, ok := l.lang[field.Key]; ok {
			fields[i].Key = v
		}
	}
	l.Logger = l.Logger.With(fields...)
	return &l
}

func (l *Logger) formatLang(msg *string, fields []zapcore.Field) {
	if l.lang != nil {
		if v, ok := l.lang[*msg]; ok {
			*msg = v
		}
		for i, field := range fields {
			if v, ok := l.lang[field.Key]; ok {
				fields[i].Key = v
			}
		}
	}
}

func (l *Logger) Trace(msg string, fields ...zap.Field) {
	l.formatLang(&msg, fields)
	l.Logger.Debug(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.formatLang(&msg, fields)
	l.Logger.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.formatLang(&msg, fields)
	l.Logger.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.formatLang(&msg, fields)
	l.Logger.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.formatLang(&msg, fields)
	l.Logger.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.formatLang(&msg, fields)
	l.Logger.Fatal(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.formatLang(&msg, fields)
	l.Logger.Panic(msg, fields...)
}
