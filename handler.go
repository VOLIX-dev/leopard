package leopard

import (
	"net/http"
)

func (l LeopardApp) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	l.router.ServeHTTP(writer, request)
}
