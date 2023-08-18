package base

import (
	"fmt"

	"github.com/baowk/dilu-core/core"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

func FmtReqId(reqId string) string {
	return fmt.Sprintf("REQID:%s", reqId)
}

func GetAcceptLanguage(c *gin.Context) string {
	return c.GetHeader("Accept-Language")
}

func GetMsgByCode(c *gin.Context, code int) string {
	if core.Cfg.Server.I18n {
		acceptLanguate := GetAcceptLanguage(c)
		tags, _, _ := language.ParseAcceptLanguage(acceptLanguate)
		if len(tags) > 0 {
			return core.I18n.GetMsg(code, tags[0].String())
		}
	}
	return core.I18n.GetMsg(code, core.Cfg.Server.GetLang())
}
