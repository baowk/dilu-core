package i18n

import "github.com/gin-gonic/gin"

var Lang ILang

type ILang interface {
	GetMsg(code int, c *gin.Context) string
	Enable() bool
	DefLang() string
}

func Register(i ILang) {
	Lang = i
}
