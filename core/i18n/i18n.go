package i18n

type I18n interface {
	GetMsg(code int, lang string) (str string)
}

func Register(i I18n) {

}
