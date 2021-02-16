package loggerx

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func GetLoggerWithCtx(c *gin.Context) *logrus.Entry {
	log, exist := c.Get("logger")
	if !exist {
		logrus.Fatal("Logger retrival error!")
	}
	return log.(*logrus.Entry)
}

// used for debugging,
func GetLoggerWithCtxAddBody(c *gin.Context, body string) *logrus.Entry {
	log, exist := c.Get("logger")
	if !exist {
		logrus.Fatal("Logger retrival error!")
	}
	lg := log.(*logrus.Entry)
	return lg.WithField("body", body)
}

func UseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request setup custom response writer to capture the response
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		if req := GetReqId(c); req == "" {
			logrus.Fatal("No request id associated with request!")
		} else {
			log := logrus.WithFields(
				logrus.Fields{
					"request-id": req,
					"route":      c.FullPath(),
				},
			)
			c.Set("logger", log)

		}
		c.Next()

		status := c.Writer.Status()
		if status != 200 {
			log := GetLoggerWithCtxAddBody(c, w.body.String())
			log.Warn(fmt.Sprintf("Request fail with status code %d", status))

		}
	}
}
