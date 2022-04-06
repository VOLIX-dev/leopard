package leopard

import "net/http"

type SimpleServer struct {
	*http.Server
}

func NewSimpleServer() *SimpleServer {
	return &SimpleServer{
		Server: &http.Server{
			Addr: ":8080",
		},
	}
}
