package zap

import (
	"github.com/132982317/profstik/pkg/utils/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2" // 日志切割归档库
	"sync"
)

// 从配置文件中获取日志的输出路径
var (
	config     = viper.Init("log")                //初始化viper并导入设置
	infoPath   = config.Viper.GetString("log")    //INFO&DEBUG&WARN级别的日志输出位置
	errorPath  = config.Viper.GetString("errorf") //ERROR和FATAL级别的日志输出位置
	LoggerPool = sync.Pool{New: func() interface{} {
		return InitLogger()
	}}
)

// InitLogger 初始化zap日志库
func InitLogger() *zap.SugaredLogger {
	//规定日志级别
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})

	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})

	//各级别通用的encoder
	encoder := getEncoder()

	//INFO级别的日志
	infoSyncer := getInfoWriter()
	infoCore := zapcore.NewCore(encoder, infoSyncer, lowPriority)

	//ERROR级别的日志
	errorSyncer := getErrorWriter()
	errorCore := zapcore.NewCore(encoder, errorSyncer, highPriority)

	//合并各种级别的日志
	core := zapcore.NewTee(infoCore, errorCore)
	logger := zap.New(core, zap.AddCaller())
	sugarLogger := logger.Sugar()
	return sugarLogger
}

// 自定义日志输出格式
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig() //生产环境下的配置
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 日志级别大写显示
	return zapcore.NewConsoleEncoder(encoderConfig)         // 输出到控制台
}

// 获取INFO的Writer
func getInfoWriter() zapcore.WriteSyncer {
	//lumberJack:日志切割归档
	lumberJackLogger := &lumberjack.Logger{
		Filename:   infoPath, // 日志输出路径
		MaxSize:    100,      // 单个日志文件的最大大小，以MB为单位
		MaxBackups: 5,        // 最多保留的备份文件的个数
		MaxAge:     30,       // 保留日志的最大天数
		Compress:   false,    // 是否压缩备份文件
	}
	return zapcore.AddSync(lumberJackLogger)
}

// 获取ERROR的Writer
func getErrorWriter() zapcore.WriteSyncer {
	//lumberJack:日志切割归档
	lumberJackLogger := &lumberjack.Logger{
		Filename:   errorPath, // 日志输出路径
		MaxSize:    100,       // 单个日志文件的最大大小，以MB为单位
		MaxBackups: 5,         // 最多保留的备份文件的个数
		MaxAge:     30,        // 保留日志的最大
		Compress:   false,     // 是否压缩备份文件
	}
	return zapcore.AddSync(lumberJackLogger)
}
