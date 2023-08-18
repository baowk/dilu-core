package core

import "net/http"

type IServer interface {
	GetEngine() http.Handler
}
