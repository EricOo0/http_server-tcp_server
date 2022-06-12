package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	_ "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var DefaultLogger *zap.Logger

func BootLogger() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder //指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	//日志级别
	errorLog := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //error级别
		return lev >= zap.ErrorLevel
	})
	infoLog := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //info和debug级别,debug级别是最低的
		return lev == zap.InfoLevel
	})
	debugLog := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //info和debug级别,debug级别是最低的
		return lev == zap.DebugLevel
	})

	// lumberjack 日志切割
	infoFileSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/info.log", //日志文件存放目录
		MaxSize:    1,                //文件大小限制,单位MB
		MaxBackups: 5,                //最大保留日志文件数量
		MaxAge:     30,               //日志文件保留天数
		Compress:   false,
	})
	infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileSyncer, zapcore.AddSync(os.Stdout)), infoLog)

	debugFileSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/debug.log", //日志文件存放目录
		MaxSize:    1,                 //文件大小限制,单位MB
		MaxBackups: 5,                 //最大保留日志文件数量
		MaxAge:     30,                //日志文件保留天数
		Compress:   false,
	})
	debugFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(debugFileSyncer, zapcore.AddSync(os.Stdout)), debugLog)

	errorFileSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/error.log", //日志文件存放目录
		MaxSize:    1,                 //文件大小限制,单位MB
		MaxBackups: 5,                 //最大保留日志文件数量
		MaxAge:     30,                //日志文件保留天数
		Compress:   false,
	})
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileSyncer, zapcore.AddSync(os.Stdout)), errorLog)

	var coreArr []zapcore.Core
	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, debugFileCore)
	coreArr = append(coreArr, errorFileCore)
	logger := zap.New(zapcore.NewTee(coreArr...), zap.AddCaller())
	DefaultLogger = logger
	DefaultLogger.Info("Boot logger successfully")
	return logger
}
