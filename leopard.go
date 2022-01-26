package leopard

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type LeopardApp struct {
	*Options
	router *mux.Router
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
