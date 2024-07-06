package logger

import (
	"os"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerConfig struct {
	Format   string
	LogType  string
	Priority string
}

func InitLogger(config LoggerConfig) *zap.Logger {
	var logger *zap.Logger
	var coreArr []zapcore.Core

	// Default values
	if config.Format == "" {
		config.Format = "normal"
	}
	if config.LogType == "" {
		config.LogType = "console"
	}
	if config.Priority == "" {
		config.Priority = "info"
	}

	// Get encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // time format
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // use different color for various log levels

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	if config.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// Log levels
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.ErrorLevel && lev >= zap.InfoLevel
	})
	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.InfoLevel && lev >= zap.DebugLevel
	})

	// File writers
	debugFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/debug.log",
		MaxSize:    128,
		MaxBackups: 3,
		MaxAge:     10,
		Compress:   false,
	})
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/info.log",
		MaxSize:    128,
		MaxBackups: 3,
		MaxAge:     10,
		Compress:   false,
	})
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/error.log",
		MaxSize:    128,
		MaxBackups: 5,
		MaxAge:     10,
		Compress:   false,
	})

	// Cores
	debugFileCore := zapcore.NewCore(encoder, os.Stdout, debugPriority)
	infoFileCore := zapcore.NewCore(encoder, os.Stdout, infoPriority)
	errorFileCore := zapcore.NewCore(encoder, os.Stdout, errorPriority)

	if config.LogType == "file" {
		debugFileCore = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(debugFileWriteSyncer, zapcore.AddSync(os.Stdout)), debugPriority)
		infoFileCore = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), infoPriority)
		errorFileCore = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), errorPriority)
	}

	switch config.Priority {
	case "info":
		coreArr = append(coreArr, infoFileCore)
		coreArr = append(coreArr, errorFileCore)
	case "error":
		coreArr = append(coreArr, errorFileCore)
	case "debug":
		coreArr = append(coreArr, debugFileCore)
		coreArr = append(coreArr, infoFileCore)
		coreArr = append(coreArr, errorFileCore)
	default:
		coreArr = append(coreArr, infoFileCore)
		coreArr = append(coreArr, errorFileCore)
	}

	logger = zap.New(zapcore.NewTee(coreArr...), zap.AddCaller()) //zap.AddCaller() is to show the line number
	return logger
}

