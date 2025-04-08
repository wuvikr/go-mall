package middleware

import (
	"bytes"
	"go-mall/common/logger"
	"go-mall/common/util"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func StartTrace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.Request.Header.Get("traceid")
		pSpanId := c.Request.Header.Get("pspanid")
		spanId := util.GenerateSpanID(c.Request.RemoteAddr)

		// 如果traceId 为空，证明是链路的发端，把它设置成此次的spanId，发端的spanId是root spanId
		if traceId == "" {
			traceId = spanId // trace 标识整个请求的链路, span则标识链路中的不同服务
		}

		c.Set("traceid", traceId)
		c.Set("pspanid", pSpanId)
		c.Set("spanid", spanId)

		c.Next()
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// 包装一下 gin.ResponseWriter，通过这种方式拦截写响应
// 让gin写响应的时候先写到 bodyLogWriter 再写gin.ResponseWriter ，
// 这样利用中间件里输出访问日志时就能拿到响应了
// https://stackoverflow.com/questions/38501325/how-to-log-response-body-in-gin
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LogAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 保存自己的body
		reqBody := []byte{}
		if c.Request.Body != nil {
			reqBody, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewReader(reqBody))
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 开始计时
		start := time.Now()
		accessLog(c, "access_start", time.Since(start), reqBody, nil)
		defer func() {
			accessLog(c, "access_end", time.Since(start), reqBody, blw.body.String())
		}()

		c.Next()

	}
}

func accessLog(c *gin.Context, accessType string, dur time.Duration, body []byte, dataOut any) {
	req := c.Request
	bodyStr := string(body)
	query := req.URL.RawQuery
	path := req.URL.Path
	// TODO: 实现 Token认证后，加上 token信息
	//
	// token := c.GetHeader("Authorization")
	logger.New(c).Info(
		"AccessLog",
		"type", accessType,
		"ip", c.ClientIP(),
		"method", req.Method,
		"path", path,
		"query", query,
		"body", bodyStr,
		"output", dataOut,
		"time(ms)", int64(dur.Milliseconds()),
	)

}

// Recovery 自定义gin recover输出
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.New(c).Error("http request broken pipe", "path", c.Request.URL.Path, "error", err, "request", string(httpRequest))
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				logger.New(c).Error("http_request_panic", "path", c.Request.URL.Path, "error", err, "request", string(httpRequest), "stack", string(debug.Stack()))

				c.AbortWithError(http.StatusInternalServerError, err.(error))
			}
		}()
		c.Next()
	}
}
