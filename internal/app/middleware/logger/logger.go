package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

const (
	TimestampFormat = "2006/01/02 15:04:05.000"
	maxAge          = 30 * 24 * time.Hour
	rotationTime    = 24 * time.Hour
	appNameKey      = "RUNTIME_APP_NAME"
)

var (
	Log *slog.Logger

	logFilePath = "/tmp"

	loggerOnce sync.Once
	loggerErr  error
)

func GetAppName() string {
	return os.Getenv(appNameKey)
}

// Logger 初始化并返回全局日志实例。
// 日志按天轮转，保留 30 天，同时输出到控制台和文件。
func Logger() (*slog.Logger, error) {
	loggerOnce.Do(func() {
		appName := GetAppName()
		if appName == "" {
			appName = "app"
		}

		// 避免环境变量意外包含目录路径。
		appName = filepath.Base(appName)

		appLogFilePath := filepath.Join(logFilePath, appName)

		if err := os.MkdirAll(logFilePath, 0755); err != nil {
			loggerErr = fmt.Errorf("create log directory: %w", err)
			return
		}

		appLogWriter, err := rotatelogs.New(
			appLogFilePath+".%Y-%m-%d.%H%M.log",
			rotatelogs.WithLocation(time.Local),
			rotatelogs.WithLinkName(appLogFilePath+".log"),
			rotatelogs.WithMaxAge(maxAge),
			rotatelogs.WithRotationTime(rotationTime),
		)
		if err != nil {
			loggerErr = fmt.Errorf("initialize rotating log writer: %w", err)
			return
		}

		handler := slog.NewJSONHandler(
			io.MultiWriter(os.Stdout, appLogWriter),
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
				ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
					if attr.Key == slog.TimeKey {
						attr.Value = slog.StringValue(
							attr.Value.Time().Format(TimestampFormat),
						)
					}
					return attr
				},
			},
		)

		Log = slog.New(handler)
	})

	return Log, loggerErr
}

// LogMiddleware 记录 Gin HTTP 请求日志。
func LogMiddleware() gin.HandlerFunc {
	logger, err := Logger()
	if err != nil {
		panic(fmt.Errorf("initialize logger middleware: %w", err))
	}

	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		latency := time.Since(startTime)

		requestID := c.Writer.Header().Get("X-Request-Id")
		traceID := c.Writer.Header().Get("TraceId")

		// 兼容请求 ID 由上游通过请求头传入的情况。
		if requestID == "" {
			requestID = c.Request.Header.Get("X-Request-Id")
		}
		if traceID == "" {
			traceID = c.Request.Header.Get("TraceId")
		}

		logger.InfoContext(
			c.Request.Context(),
			"http request",
			slog.Int("status_code", c.Writer.Status()),
			slog.String("x-request-id", requestID),
			slog.String("trace_id", traceID),
			slog.String("latency_time", latency.String()),
			slog.Float64("latency_ms", float64(latency)/float64(time.Millisecond)),
			slog.String("client_ip", c.ClientIP()),
			slog.String("req_method", c.Request.Method),
			slog.String("req_uri", c.Request.RequestURI),
		)
	}
}
