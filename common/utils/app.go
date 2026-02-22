package utils

import (
	"github.com/baowk/dilu-core/common/consts"
	"github.com/gin-gonic/gin"
)

func GetReqId(c *gin.Context) string {
	reqId := c.GetString(consts.REQ_ID)
	if reqId == "" {
		reqId = GenString()
		c.Set(consts.REQ_ID, reqId)
	}
	return reqId
}
