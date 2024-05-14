package midwares

import (
	"QA-System/app/apiException"

	"net/http"
	"os"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ErrHandler() gin.HandlerFunc {
	logFilePath := "app.log"

	// 检查日志文件是否存在
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		// 如果日志文件不存在，则创建新的日志文件
		_, err := os.Create(logFilePath)
		if err != nil {
			// 创建日志文件失败，记录错误并返回空的中间件处理函数
			zap.S().Error("Failed to create log file:", err)
			return func(c *gin.Context) {}
		}
	}

	// 打开日志文件
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		// 打开日志文件失败，记录错误并返回空的中间件处理函数
		zap.S().Error("Failed to open log file:", err)
		return func(c *gin.Context) {}
	}
	writeSyncer := zapcore.AddSync(logFile)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync()

	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Handle panic and log the error
				stack := debug.Stack()
				logger.Error("Panic recovered",
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.Any("panic", r),
					zap.ByteString("stacktrace", stack),
				)

				c.JSON(http.StatusInternalServerError, gin.H{
					"code":  http.StatusInternalServerError,
					"msg": apiException.ServerError.Msg,
				})
				c.Abort()
			}
		}()

		c.Next()
		if length := len(c.Errors); length > 0 {
			e := c.Errors[length-1]
			err := e.Err
			if err != nil {
				// TODO 建立日志系统
				var logLevel zapcore.Level
				switch e.Type {
				case gin.ErrorTypePublic:
					logLevel = zapcore.ErrorLevel
				case gin.ErrorTypeBind:
					logLevel = zapcore.WarnLevel
				case gin.ErrorTypePrivate:
					logLevel = zapcore.DebugLevel
				default:
					logLevel = zapcore.InfoLevel
				}
				logger.Check(logLevel, "Error reported").Write(
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.Error(err),
				)
				return
			}
		}
	}
}


// HandleNotFound
//
//	404处理
func HandleNotFound(c *gin.Context) {
	err := apiException.NotFound
	c.JSON(err.StatusCode, err)
}
