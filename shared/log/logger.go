package log

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *logrus.Logger

func InitLogger(name, logPath, logLevel string) {
	Logger = logrus.New()

	// 设置日志格式
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true, // 如果你想在控制台输出颜色，可以设置为true
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// 获取包名
			pkg := path.Base(path.Dir(f.File))

			filename := path.Base(f.File)
			return "", fmt.Sprintf(" %s/%s:%d", pkg, filename, f.Line)
		},
	})

	// 设置报告调用者信息
	Logger.SetReportCaller(true)

	// 设置日志输出到 lumberjack，实现日志轮转和压缩
	lj := &lumberjack.Logger{
		Filename:   path.Join(logPath, name), // 日志文件路径
		MaxSize:    50,                       // 每个日志文件的最大大小，单位MB
		MaxBackups: 3,                        // 保留旧文件的最大个数
		MaxAge:     7,                        // 保留旧文件的最大天数
		Compress:   true,                     // 是否压缩/归档旧文件
		LocalTime:  true,
	}

	// 同时输出到控制台和文件
	Logger.SetOutput(io.MultiWriter(os.Stdout, lj))

	// 设置日志级别
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel // 默认级别
	}
	Logger.SetLevel(level)

	Logger.Info("Logger initialized successfully")
}

func GetLogger() *logrus.Logger {
	return Logger
}

func WithFields(fields logrus.Fields) *logrus.Logger {
	Logger.WithFields(fields)
	return Logger
}

func Trace(f interface{}, v ...interface{}) {
	if Logger == nil {
		logrus.Trace(formatLog(f, v...))
		return
	}
	Logger.Trace(formatLog(f, v...))
}

func formatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return fmt.Sprintf(msg, v...)
}
