package leopard

import "net/http"

type SimpleServer struct {
	*http.Server
}

func NewSimpleServer(addr string) *SimpleServer {
	return &SimpleServer{
		Server: &http.Server{
			Addr: addr,
		},
	}
}
