package leopard

import "net/http"

type SimpleServer struct {
	*http.Server
}

func (s *SimpleServer) Start(h http.Handler) error {
	s.Handler = h
	return s.ListenAndServe()
}
