package leopard

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

type MiddlewareFunc mux.MiddlewareFunc

// GET handler register
func (a *LeopardApp) GET(p string, name string, h func(r *Context), middleware ...MiddlewareFunc) {
	a.AddRoute(http.MethodGet, p, name, h)
}

// POST register a route with the method POST
func (a *LeopardApp) POST(p string, name string, h func(r *Context), middleware ...MiddlewareFunc) {
	a.AddRoute(http.MethodPost, p, name, h)
}

// PUT register a route with the method PUT
func (a *LeopardApp) PUT(p string, name string, h func(r *Context), middleware ...MiddlewareFunc) {
	a.AddRoute(http.MethodPut, p, name, h)
}

// DELETE register a route with the method DELETE
func (a *LeopardApp) DELETE(p string, name string, h func(r *Context), middleware ...MiddlewareFunc) {
	a.AddRoute(http.MethodDelete, p, name, h)
}

// PATCH register a reoute with the method PATCH
func (a *LeopardApp) PATCH(p string, name string, h func(r *Context), middleware ...MiddlewareFunc) {
	a.AddRoute(http.MethodPatch, p, name, h)
}

// AddRoute adds a route to the route manager
// This is mainly called by methods as GET, POST, PUT, DELETE and PATCH
// However if needed a user could register a custom method name (or one we did not include)
func (a *LeopardApp) AddRoute(method string, p string, name string, h func(r *Context), middleware ...MiddlewareFunc) {
	r := a.router.NewRoute()

	r.Methods(method)
	r.Path(p)
	r.Name(name)
	r.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		context := NewContext(w, r, a)

		defer func() {
			//if r := recover(); r != nil {
			//	err := context.Error(fmt.Errorf("%v", r))
			//
			//	if err != nil {
			//		return
			//	}
			//}
		}()
		h(context)
	})
}

// StaticDir register a static directory
func (a *LeopardApp) StaticDir(p string, root http.FileSystem) {
	pa := path.Join(p, a.Prefix)

	h := stripAsset(pa, a.fileServer(root))
	a.router.PathPrefix(pa).Handler(h)
}

// fileServer creates a file server and returns its handler
func (a *LeopardApp) fileServer(fs http.FileSystem) http.Handler {
	fsh := http.FileServer(fs)
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			w.WriteHeader(404)
			return
		}

		stat, _ := f.Stat()
		maxAge := "31536000"
		w.Header().Add("ETag", fmt.Sprintf("%x", stat.ModTime().UnixNano()))
		w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%s", maxAge))
		fsh.ServeHTTP(w, r)
	})
	//
	//if a.CompressFiles {
	//	return handlers.CompressHandler(baseHandler)
	//}

	return baseHandler
}

// stripAsset strips path of assets on a file server
func stripAsset(path string, handler http.Handler) http.Handler {
	if path == "" {
		return handler
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		up := r.URL.Path
		up = strings.TrimPrefix(up, path)
		up = strings.TrimSuffix(up, "/")

		u, err := url.Parse(up)
		if err != nil {
			w.WriteHeader(404)
			return
		}

		r.URL = u
		handler.ServeHTTP(w, r)
	})
}
