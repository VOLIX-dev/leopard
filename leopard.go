package leopard

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	cacheDrivers "github.com/volix-dev/leopard/caching/drivers"
	"github.com/volix-dev/leopard/files"
	"github.com/volix-dev/leopard/templating"
	"github.com/volix-dev/leopard/templating/drivers"
	"net/http"
	"strconv"
)

type LeopardApp struct {
	*Options
	router   *mux.Router
	server   *SimpleServer
	settings map[string]SettingValue

	TemplateDriver drivers.TemplatingDriver
	Cache          *Caching
	FileDriver     files.Driver

	ContextCreator func(r *http.Request, w http.ResponseWriter, a *LeopardApp) ContextInterface
}

func New() (*LeopardApp, error) {
	app := &LeopardApp{
		router:         mux.NewRouter(),
		TemplateDriver: templating.TwigCreator(),
		ContextCreator: func(r *http.Request, w http.ResponseWriter, a *LeopardApp) ContextInterface {
			return &Context{
				request:        r,
				responseWriter: w,
				a:              a,
			}
		},
	}

	err := godotenv.Load()

	err = app.TemplateDriver.Load(
		EnvSettingD("TEMPLATE_PATH", "./templates").GetValue().(string),
		app.router,
	)
	if err != nil {
		return nil, err
	}

	driverName := EnvSettingD("CACHE_DRIVER", "memory").GetValue().(string)

	switch driverName {
	case "redis":
		port, err := strconv.Atoi(EnvSettingD("REDIS_PORT", "6379").GetValue().(string))
		db, err := strconv.Atoi(EnvSettingD("REDIS_DB", "0").GetValue().(string))

		if err != nil {
			return nil, err
		}

		cache, err := newCaching(driverName, cacheDrivers.RedisSettings{
			Host:     EnvSettingD("REDIS_HOST", "localhost").GetValue().(string),
			Port:     port,
			Password: EnvSettingD("REDIS_PASSWORD", "").GetValue().(string),
			Database: db,
		})

		if err != nil {
			return nil, err
		}

		app.Cache = cache

		break

	case "memory":
		cache, err := newCaching(driverName, nil)

		if err != nil {
			return nil, err
		}

		app.Cache = cache

		break

	default:
		return nil, errors.New("cache driver not found")
	}

	err = app.Cache.open()
	if err != nil {
		return nil, err
	}

	fileDriver, err := getFileDriver()
	if err != nil {
		return nil, err
	}
	app.FileDriver = fileDriver

	return app, nil
}

// Serve starts the server
// It will listen on the port specified in the options.
// Should only be called once.
func (a *LeopardApp) Serve() error {
	a.server = NewSimpleServer(EnvSettingD("LISTEN", ":8080").GetValue().(string))
	a.server.Handler = a.router
	err := a.server.ListenAndServe()

	return err
}

// GetRouter gets the mux router
func (a *LeopardApp) GetRouter() *mux.Router {
	return a.router
}

// Close closes all connections and resources.
// It is not safe to use a leopard app after it has been closed.
// Close should be called before exiting the program.
//
// Close gets called automatically when there is an interrupt signal.
// If you have any other way of closing the app, you should call this function.
func (a *LeopardApp) Close() error {
	return a.server.Close()
}
