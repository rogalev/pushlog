package logging

import (
	"github.com/rogalev/pushlog/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var logger *zap.Logger
var once sync.Once
var logFilePath string

func GetInstance() *zap.Logger {
	once.Do(func() {
		logger = createLogger()
	})
	return logger
}

func SetupConfig(cfg config.Config) {
	logFilePath = cfg.LogFile
}

func createLogger() *zap.Logger {

	c := zap.NewProductionEncoderConfig()
	c.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(c)
	consoleEncoder := zapcore.NewConsoleEncoder(c)

	logFile, _ := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)

	defaultLogLevel := zapcore.DebugLevel

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}
