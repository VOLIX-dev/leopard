package leopard

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

// GET handler register
func (l *LeopardApp) GET(p string, h func(r *http.Request)) {
	l.AddRoute(http.MethodGet, p, h)
}

// POST register a route with the method POST
func (l *LeopardApp) POST(p string, h func(r *http.Request)) {
	l.AddRoute(http.MethodPost, p, h)
}

// PUT register a route with the method PUT
func (l *LeopardApp) PUT(p string, h func(r *http.Request)) {
	l.AddRoute(http.MethodPut, p, h)
}

// DELETE register a route with the method DELETE
func (l *LeopardApp) DELETE(p string, h func(r *http.Request)) {
	l.AddRoute(http.MethodDelete, p, h)
}

// PATCH register a reoute with the method PATCH
func (l *LeopardApp) PATCH(p string, h func(r *http.Request)) {
	l.AddRoute(http.MethodPatch, p, h)
}

// AddRoute adds a route to the route manager
// This is mainly called by methods as GET, POST, PUT, DELETE and PATCH
// However if needed a user could register a custom method name (or one we did not include)
func (l *LeopardApp) AddRoute(method string, p string, h func(r *http.Request)) {
	r := l.router.NewRoute()

	r.Methods(method)
	r.Path(p)
	r.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(r)
	})
}

// StaticDir register a static directory
func (l *LeopardApp) StaticDir(p string, root http.FileSystem) {
	pa := path.Join(p, l.Prefix)

	h := stripAsset(pa, l.fileServer(root), l)
	l.router.PathPrefix(pa).Handler(h)
}

// fileServer creates a file server and returns its handler
func (l *LeopardApp) fileServer(fs http.FileSystem) http.Handler {
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
	//if l.CompressFiles {
	//	return handlers.CompressHandler(baseHandler)
	//}

	return baseHandler
}

// stripAsset strips path of assets on a file server
func stripAsset(path string, handler http.Handler, l *LeopardApp) http.Handler {
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
