package base

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func FmtReqId(reqId string) string {
	return fmt.Sprintf("REQID:%s", reqId)
}

func GetAcceptLanguage(c *gin.Context) string {
	return c.GetHeader("Accept-Language")
}

// func GetMsgByCode(c *gin.Context, code int) string {

// }
