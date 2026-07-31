package utils

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Initzap() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("打开日志文件失败" + err.Error())
	}
	wrireSyncFile := zapcore.AddSync(logFile)
	writeSyncConsole := zapcore.AddSync(os.Stdout)
	encoderconfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "msg",
		EncodeTime:     timeformat,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
	encoder := zapcore.NewJSONEncoder(encoderconfig)
	level := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, wrireSyncFile, level),
		zapcore.NewCore(encoder, writeSyncConsole, level))
	Log = zap.New(core, zap.AddCaller())
}
func timeformat(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
