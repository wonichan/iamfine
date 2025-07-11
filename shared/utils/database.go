package utils

import (
	"fmt"
	"hupu/shared/config"
	"hupu/shared/models"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func CloseDB() error {
	db, err := instance.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func InitDB() error {
	// 设置默认时区为北京时间
	loc, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = loc
	cfg := config.GlobalConfig.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FShanghai",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	// 配置GORM日志 - 使用项目自定义日志系统
	// 从配置文件读取日志级别，如果未配置则默认为 info
	logLevel := cfg.LogLevel
	if logLevel == "" {
		logLevel = "info"
	}

	gormLogger := NewGormLogger().
		SetSlowThreshold(time.Second).
		SetIgnoreRecordNotFoundError(true).
		SetParameterizedQueries(false).
		LogMode(ParseLogLevel(logLevel))

	instance, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().In(loc)
		},
		Logger: gormLogger,
	})
	if err != nil {
		return err
	}

	// 自动迁移数据库表
	err = instance.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
		&models.Like{},
		&models.Follow{},
		&models.Notification{},
		&models.PostFavorite{},
		&models.PostRating{},
		&models.Topic{},
	)
	if err != nil {
		return err
	}

	return nil
}

var instance *gorm.DB

var onceDB sync.Once

func GetDB() *gorm.DB {
	onceDB.Do(func() {
		if err := InitDB(); err != nil {
			panic(err)
		}
	})
	return instance
}
