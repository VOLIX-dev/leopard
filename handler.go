package leopard

import (
	"net/http"
)

func (a LeopardApp) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	a.router.ServeHTTP(writer, request)
}
