package i18n

var Lang ILang

type ILang interface {
	GetMsg(code int, lang string) (str string)
}

func Register(i ILang) {
	Lang = i
}
