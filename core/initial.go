package core

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"poem-bot/global"
	"poem-bot/util"
	"time"
)

func Run() {
	Init()
}

func Init() {
	global.VP = Viper()
	global.LOG = Zap()
	global.DB = Gorm()

	zap.ReplaceGlobals(global.LOG)
	defer func(LOG *zap.Logger) {
		_ = LOG.Sync()
	}(global.LOG)
}

func Viper() *viper.Viper {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	return v
}

func Zap() *zap.Logger {
	logDir := util.GetOrDefault("zap.directory", "log")
	if exists, _ := util.PathExists(logDir); !exists {
		fmt.Printf("create %v directory\n", logDir)
		_ = os.Mkdir(logDir, os.ModePerm)
	}

	// 日志时间格式
	encodeTime := func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	// 日志格式配置
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "name",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     encodeTime,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		RotateLogFile.GetWriteSyncer(logDir),
		zap.InfoLevel,
	)
	logger := zap.New(core, zap.AddCaller())
	return logger
}

func Gorm() *gorm.DB {
	mc := global.VP.Sub("mysql")
	if mc == nil {
		global.LOG.Error("mysql config is nil")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		global.VP.Get("mysql.username"),
		global.VP.Get("mysql.password"),
		global.VP.Get("mysql.host"),
		global.VP.Get("mysql.port"),
		global.VP.Get("mysql.db-name"),
		global.VP.Get("mysql.config"),
	)

	m := mysql.New(mysql.Config{DSN: dsn})
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // 慢 SQL 阈值
			LogLevel:      logger.Silent, // 设置日志级别为Silent
			Colorful:      false,         // 禁用彩色打印
		},
	)
	db, err := gorm.Open(m, &gorm.Config{Logger: newLogger})
	if err != nil {
		panic(err)
	}
	db.InstanceSet("gorm:table_options", "ENGINE=InnoDB")
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(global.VP.GetInt("mysql.max-idle-conn"))
	sqlDB.SetMaxOpenConns(global.VP.GetInt("mysql.max-open-conn"))

	global.LOG.Info("Database connection successful")
	return db
}
