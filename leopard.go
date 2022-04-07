package leopard

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"leopard/templating"
	"leopard/templating/drivers"
)

type LeopardApp struct {
	*Options
	router   *mux.Router
	server   *SimpleServer
	settings map[string]SettingValue

	TemplateDriver drivers.TemplatingDriver
}

func New() (*LeopardApp, error) {
	app := &LeopardApp{
		router:         mux.NewRouter(),
		TemplateDriver: templating.TwigCreator(),
	}

	err := godotenv.Load()

	err = app.TemplateDriver.Load(
		EnvSettingD("TEMPLATE_PATH", "./templates").GetValue().(string),
		app.router,
	)
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

func (a *LeopardApp) GetRouter() *mux.Router {
	return a.router
}
