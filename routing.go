package leopard

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
)

type MiddlewareFunc func(context ContextInterface)

// GET handler register
func (a *LeopardApp) GET(p string, h func(r ContextInterface), extras ...any) {
	a.AddRoute(http.MethodGet, p, h, extras...)
}

// POST register a route with the method POST
func (a *LeopardApp) POST(p string, h func(r ContextInterface), extras ...any) {
	a.AddRoute(http.MethodPost, p, h, extras...)
}

// PUT register a route with the method PUT
func (a *LeopardApp) PUT(p string, h func(r ContextInterface), extras ...any) {
	a.AddRoute(http.MethodPut, p, h, extras)
}

// DELETE register a route with the method DELETE
func (a *LeopardApp) DELETE(p string, h func(r ContextInterface), extras ...any) {
	a.AddRoute(http.MethodDelete, p, h, extras)
}

// PATCH register a reoute with the method PATCH
func (a *LeopardApp) PATCH(p string, h func(r ContextInterface), extras ...any) {
	a.AddRoute(http.MethodPatch, p, h, extras)
}

func (a *LeopardApp) Group(p string, groupHandler func(group RouteGroup), extras ...any) RouteGroup {
	name, middleware := parseExtras(extras)
	group := RouteGroup{
		prefix:     p,
		namePrefix: name,
		app:        a,
		middleware: middleware,
	}
	groupHandler(group)

	return group
}

// AddRoute adds a route to the route manager
// This is mainly called by methods as GET, POST, PUT, DELETE and PATCH
// However if needed a user could register a custom method name (or one we did not include)
func (a *LeopardApp) AddRoute(method string, p string, h func(r ContextInterface), extras ...any) {
	name, middleware := parseExtras(extras)

	a.addRoute(method, p, h, name, middleware)
}

func (a *LeopardApp) addRoute(method string, p string, h func(r ContextInterface), name *string, middleware []MiddlewareFunc) {
	r := a.router.NewRoute()

	r.Methods(method)
	r.Path(p)

	fmt.Println(r.GetPathTemplate())

	if name != nil {
		r.Name(*name)
	}

	r.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := a.ContextCreator(r, w, a)

		defer func() {
			if r := recover(); r != nil {
				err := context.Error(fmt.Errorf("%v", r))

				if err != nil {
					return
				}
			}
		}()

		for _, m := range middleware {
			m(context)

			if context.Aborted() {
				return
			}
		}

		h(context)
	})
}

// StaticDir register a static directory
func (a *LeopardApp) StaticDir(p string, root string) {
	h := a.fileServer(root, p)
	a.router.PathPrefix(p).Handler(h)
}

// fileServer creates a file server and returns its handler
func (a *LeopardApp) fileServer(rootDir string, p string) http.Handler {
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open(path.Join(rootDir,
			strings.ReplaceAll(
				strings.TrimPrefix(r.URL.Path, p),
				"../",
				"",
			)),
		)

		if os.IsNotExist(err) {
			w.WriteHeader(404)
			return
		}

		stat, _ := f.Stat()
		maxAge := "31536000"
		w.Header().Add("ETag", fmt.Sprintf("%x", stat.ModTime().UnixNano()))
		w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%s", maxAge))
		w.Header().Add("Content-Length", fmt.Sprintf("%d", stat.Size()))

		http.ServeContent(w, r, f.Name(), stat.ModTime(), f)
	})
	//
	//if a.CompressFiles {
	//	return handlers.CompressHandler(baseHandler)
	//}

	return baseHandler
}

type RouteGroup struct {
	prefix     string
	namePrefix *string
	middleware []MiddlewareFunc
	app        *LeopardApp
}

func (r RouteGroup) GET(p string, h func(r ContextInterface), extras ...any) {
	r.addRoute(
		http.MethodGet,
		p,
		h,
		extras...,
	)
}

// POST register a route with the method POST
func (r RouteGroup) POST(p string, h func(r ContextInterface), extras ...any) {
	r.addRoute(
		http.MethodPost,
		p,
		h,
		extras...,
	)
}

// PUT register a route with the method PUT
func (r RouteGroup) PUT(p string, h func(r ContextInterface), extras ...any) {
	r.addRoute(
		http.MethodPut,
		p,
		h,
		extras...,
	)
}

// DELETE register a route with the method DELETE
func (r RouteGroup) DELETE(p string, h func(r ContextInterface), extras ...any) {
	r.addRoute(
		http.MethodDelete,
		p,
		h,
		extras...,
	)
}

// PATCH register a route with the method PATCH
func (r RouteGroup) PATCH(p string, h func(r ContextInterface), extras ...any) {
	r.addRoute(
		http.MethodPatch,
		p,
		h,
		extras...,
	)
}

// Group creates a new RouteGroup with a prefix
func (r RouteGroup) Group(prefix string, groupHandler func(group RouteGroup), extras ...any) RouteGroup {
	name, middleware := parseExtras(extras)

	group := RouteGroup{
		prefix:     path.Join(r.prefix, prefix),
		namePrefix: r.addNamePrefix(name),
		app:        r.app,
		middleware: middleware,
	}
	groupHandler(group)

	return group
}

func (r RouteGroup) addRoute(method string, p string, h func(r ContextInterface), extras ...any) {
	name, middleware := parseExtras(extras)

	r.app.addRoute(
		method,
		path.Join(r.prefix, p),
		h,
		r.addNamePrefix(name),
		append(r.middleware, middleware...),
	)
}

func (r RouteGroup) addNamePrefix(name *string) *string {
	if name == nil {
		return nil
	}

	if r.namePrefix == nil {
		return name
	}

	temp := *r.namePrefix + "." + *name
	return &temp
}

func parseExtras(extras []any) (name *string, middleware []MiddlewareFunc) {
	for _, e := range extras {
		switch e.(type) {
		case string:
			temp := e.(string)
			name = &temp
			break

		case MiddlewareFunc:
			middleware = append(middleware, e.(MiddlewareFunc))
			break

		case []MiddlewareFunc:
			middleware = append(middleware, e.([]MiddlewareFunc)...)
			break
		}
	}
	return
}
