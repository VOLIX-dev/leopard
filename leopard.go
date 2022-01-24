package leopard

import (
	"net/http"
)

type LeopardApp struct {
}

func (l LeopardApp) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

}
