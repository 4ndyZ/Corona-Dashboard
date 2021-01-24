package main

import (
	"os"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger struct to hold refs
type Logger struct {
	Logger *zap.SugaredLogger
}

func (l *Logger) Initialize(d bool, f string) {
	pe := zap.NewProductionEncoderConfig()
	// Set time format to ISO8601
	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	// Create console and log file encoder
	fileEncoder := zapcore.NewJSONEncoder(pe)
	consoleEncoder := zapcore.NewConsoleEncoder(pe)
	// Check if debugging mode is enabled
	level := zap.InfoLevel
	if d {
		level = zap.DebugLevel
	}
	// Combine console and file encoder with the log outputs to the core logger
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(&lumberjack.Logger{
		  Filename:   f,
		  MaxSize:    500, // Megabytes
		  MaxBackups: 3,
		  MaxAge:     28, // Days
		}), level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)
	// Create logger from core
	logger := zap.New(core)
	// Set logger to struct
	l.Logger = logger.Sugar()
}
