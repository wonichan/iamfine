package utils

import (
	"context"
	"errors"
	"hupu/shared/log"
	"strings"
	"time"

	"gorm.io/gorm/logger"
)

// GormLogger 实现 gorm.io/gorm/logger.Interface 接口
// 将 GORM 日志重定向到项目的 logrus 日志系统
type GormLogger struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
	LogLevel                  logger.LogLevel
}

// NewGormLogger 创建新的 GORM 日志适配器
func NewGormLogger() *GormLogger {
	return &GormLogger{
		SlowThreshold:             time.Second,
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      false,
		LogLevel:                  logger.Info,
	}
}

// LogMode 设置日志级别
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 输出信息级别日志
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		log.GetLogger().Infof(msg, data...)
	}
}

// Warn 输出警告级别日志
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		log.GetLogger().Warnf(msg, data...)
	}
}

// Error 输出错误级别日志
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		log.GetLogger().Errorf(msg, data...)
	}
}

// Trace 输出 SQL 执行日志
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 处理参数化查询
	if l.ParameterizedQueries {
		sql = logger.ExplainSQL(sql, nil, `'`, nil)
	}

	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		// 错误日志
		log.GetLogger().WithField("elapsed", elapsed).WithField("rows", rows).Errorf("[GORM] %s [%v]\n%s", err, elapsed, sql)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		// 慢查询日志
		log.GetLogger().WithField("elapsed", elapsed).WithField("rows", rows).Warnf("[GORM] SLOW SQL >= %v [%v]\n%s", l.SlowThreshold, elapsed, sql)
	case l.LogLevel == logger.Info:
		// 普通 SQL 日志
		log.GetLogger().WithField("elapsed", elapsed).WithField("rows", rows).Infof("[GORM] [%v] %s", elapsed, sql)
	}
}

// SetSlowThreshold 设置慢查询阈值
func (l *GormLogger) SetSlowThreshold(threshold time.Duration) *GormLogger {
	l.SlowThreshold = threshold
	return l
}

// SetIgnoreRecordNotFoundError 设置是否忽略记录未找到错误
func (l *GormLogger) SetIgnoreRecordNotFoundError(ignore bool) *GormLogger {
	l.IgnoreRecordNotFoundError = ignore
	return l
}

// SetParameterizedQueries 设置是否使用参数化查询
func (l *GormLogger) SetParameterizedQueries(parameterized bool) *GormLogger {
	l.ParameterizedQueries = parameterized
	return l
}

// ParseLogLevel 解析日志级别字符串
func ParseLogLevel(level string) logger.LogLevel {
	switch strings.ToLower(level) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info // 默认为 Info 级别
	}
}
