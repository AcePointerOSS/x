package loggerx

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetLogger(c *gin.Context) *logrus.Entry {
	log, exist := c.Get("logger")
	if !exist {
		logrus.Fatal("Logger retrival error!")
	}

	return log.(*logrus.Entry)
}
func UseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if req := GetReqId(c); req == "" {
			logrus.Fatal("No request id associated with request!")
		} else {

			log := logrus.WithField("request-id", req)
			c.Set("logger", log)
		}
		c.Next()
	}
}
