package leopard

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type LeopardApp struct {
	*Options
	router   *mux.Router
	server   *SimpleServer
	settings map[string]SettingValue
}

func New() (*LeopardApp, error) {
	app := &LeopardApp{
		router: mux.NewRouter(),
	}

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *LeopardApp) Serve() {
	a.server = NewSimpleServer()
	a.server.Handler = a.router
	a.server.ListenAndServe()
}
