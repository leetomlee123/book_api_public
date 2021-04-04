package logrus

import (
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type LogDto struct {
	UserName    string `json:"userName"`
	StatusCode  string `json:"statusCode"`
	LatencyTime string `json:"latencyTime"`
	ClientIP    string `json:"clientIP"`
	ReqMethod   string `json:"reqMethod"`
	ReqUri      string `json:"reqUri"`
	Time        string `json:"time"`
}

// 日志记录到文件
func LoggerToFile() gin.HandlerFunc {

	//logFilePath := config.Log_FILE_PATH
	//logFileName := config.LOG_FILE_NAME

	//日志文件
	//fileName := path.Join(logFilePath, logFileName)

	//写入文件
	//src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	//if err != nil {
	//	fmt.Println("err", err)
	//}

	//实例化
	logger := logrus.New()

	//设置输出
	//logger.Out = src
	//logger.Out=LoggerToES

	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{})

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}

// 日志记录到 MongoDB
func LoggerToMongo() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 ES
func LoggerToES() gin.HandlerFunc {
	//实例化
	logger := logrus.New()

	//设置输出
	//logger.Out = src
	//logger.Out=LoggerToES

	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()
		//用户
		username, err := c.Get("username")
		if true != err {
			username = ""
		}
		// 日志格式
		logger.Infof("| %3d | %13v | %15s | %s | %s |%s",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
			username,
		)
		if statusCode == 200 {
			models.LogDB.Insert(LogDto{
				username.(string),
				strconv.Itoa(statusCode),
				strconv.FormatInt(latencyTime.Milliseconds(), 10),
				clientIP,
				reqMethod,
				reqUri,
				time.Now().Format("2006-01-02 15:04:05"),
			})
		}

	}
}

// 日志记录到 MQ
func LoggerToMQ() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
