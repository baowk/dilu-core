package i18n

import "github.com/gin-gonic/gin"

var Lang ILang

type ILang interface {
	GetMsg(code int, c *gin.Context) string
	//SetEnable(enable bool)
	Enable()
	//SetDefLang(lang string)
	DefLang() string
}

func Register(i ILang) {
	Lang = i
}
