package loggerx

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

func GetReqId(c *gin.Context) string {
	return c.GetHeader("request-id")
}

func UseRequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// dont preface with x as per RFC 6648 (June 2012)
		reqId := c.GetHeader("request-id")
		if reqId == "" {
			reqId = uuid.Must(uuid.NewV4()).String()
		}
		c.Header("request-id", reqId)
		c.Next()
	}
}
