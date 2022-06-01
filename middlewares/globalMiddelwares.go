package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/heejinzzz/logger"
	"strconv"
)

var logFilename = "./webLog.log"

// CrossDomain 允许跨域请求数据
func CrossDomain(c *gin.Context) {
	origin := c.Request.Header.Get("Origin") //请求头部
	if origin != "" {
		// 将该域添加到allow-origin中
		c.Header("Access-Control-Allow-Origin", origin) //
		c.Header("Access-Control-Allow-Methods",
			"POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers",
			"Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers",
			"Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, "+
				"Cache-Control, Content-Language, Content-Type")

		//允许客户端传递校验信息比如 cookie
		c.Header("Access-Control-Allow-Credentials", "true")
	}
}

// PingLog 记录用户访问网站的IP、时间、URL，写入本地的日志文件
func PingLog(c *gin.Context) {
	// 指定本地日志文件名

	urlPath := c.Request.URL.Path // 用户访问使用的url
	urlRawPath := c.Request.URL.RawQuery
	if len(urlRawPath) != 0 {
		urlPath = urlPath + "?" + urlRawPath
	}
	clientIP := c.ClientIP()        // 用户的IP
	method := c.Request.Method      // 用户的访问方式：GET/POST
	statusCode := c.Writer.Status() // HTTP响应码
	errors := c.Errors.String()     // 错误消息
	logContent := "{C-webLog}\tclientIP:" + clientIP + "\turl:" + urlPath + "\tmethod:" + method +
		"\tHTTPResponseCode:" + strconv.Itoa(statusCode)
	logger.RecordInfo(logFilename, logContent)
	if len(errors) != 0 {
		logger.RecordError(logFilename, logContent+"\terrors:"+errors)
	}
}
