package leopard

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/volix-dev/leopard/helpers"
	"github.com/volix-dev/leopard/templating/drivers"
	"io/ioutil"
	"net/http"
	"runtime/debug"
)

type ContextInterface interface {
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	App() *LeopardApp

	JsonStatus(status int, data interface{}) error
	Json(data interface{}) error
	Error(err error) error
	Status(status int)
	Ok()
	NotFound()
	Unauthorized()
	BadRequest()
	Redirect(url string)
	Write(data []byte) (int, error)
	WriteString(data string) (int, error)
	WriteStringF(format string, args ...interface{}) (int, error)
	WriteJson(data interface{}) error
	Read(data []byte) (int, error)
	ReadAll() ([]byte, error)
	ReadString() (string, error)
	ReadJson(data interface{}) error
	ReadForm() map[string][]string
	ReadFormValue(key string) string
	SetHeader(key, value string)
	SetHeaders(headers map[string][]string)
	GetHeader(key string) string
	GetHeaders() map[string][]string
	GetParam(key string) string
	HasParam(key string) bool
	GetParams() map[string]string
	GetQuery(key string) string
	Queries() map[string][]string
	GetCookie(key string) (*http.Cookie, error)
	SetCookie(key string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool)
	SetResponseCookie(cookies ...*http.Cookie)
	RenderTemplate(template string, data map[string]drivers.Value) error

	// Used for middleware only

	Aborted() bool
	Abort()
}

type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	vars           map[string]string
	a              *LeopardApp

	abort bool
}

func NewContext(w http.ResponseWriter, r *http.Request, a *LeopardApp) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
		vars:           mux.Vars(r),
		a:              a,
	}
}

// Request returns the current http.Request
func (c *Context) Request() *http.Request {
	return c.request
}

// ResponseWriter returns the http.ResponseWriter
func (c *Context) ResponseWriter() http.ResponseWriter {
	return c.responseWriter
}

// App returns the current LeopardApp
func (c *Context) App() *LeopardApp {
	return c.a
}

// JsonStatus returns the status code and the data marshaled to json.
func (c *Context) JsonStatus(status int, data interface{}) error {
	c.responseWriter.Header().Set("Content-Type", "application/json")
	c.responseWriter.WriteHeader(status)

	err := json.NewEncoder(c.responseWriter).Encode(data)

	return err
}

// Json responds with a 200 and the provided interface marshalled to json.
func (c *Context) Json(data interface{}) error {
	return c.JsonStatus(http.StatusOK, data)
}

// Error responds with a 500 and the error's message in a json.
func (c *Context) Error(err error) error {
	return c.JsonStatus(http.StatusInternalServerError, map[string]interface{}{
		"message":    err.Error(),
		"stacktrace": helpers.SerializeStack(debug.Stack()),
	})
}

// Status responds with the provided status code.
func (c *Context) Status(status int) {
	c.ResponseWriter().WriteHeader(status)
}

// Ok responds with a 200 status code.
func (c *Context) Ok() {
	c.Status(http.StatusOK)
}

// NotFound responds with a 404 status code.
func (c *Context) NotFound() {
	c.Status(http.StatusNotFound)
}

// Unauthorized responds with a 401 status code.
func (c *Context) Unauthorized() {
	c.Status(http.StatusUnauthorized)
}

// BadRequest responds with a 400 status code.
func (c *Context) BadRequest() {
	c.Status(http.StatusBadRequest)
}

// Redirect responds with a 302 status code and redirects to the provided url.
func (c *Context) Redirect(url string) {
	http.Redirect(c.ResponseWriter(), c.Request(), url, http.StatusFound)
}

// Writes

// Write writes an array of bytes directly to the response writer.
func (c *Context) Write(data []byte) (int, error) {
	return c.responseWriter.Write(data)
}

// WriteString writes a string directly to the response writer.
func (c *Context) WriteString(data string) (int, error) {
	return c.responseWriter.Write([]byte(data))
}

// WriteStringF writes a formatted string directly to the response writer.
func (c *Context) WriteStringF(format string, a ...interface{}) (int, error) {
	return c.WriteString(fmt.Sprintf(format, a...))
}

// WriteJson marshals the provided interface to json and writes it to the response writer.
func (c *Context) WriteJson(data interface{}) error {
	return json.NewEncoder(c.responseWriter).Encode(data)
}

// Reads
func (c *Context) Read(data []byte) (int, error) {
	return c.request.Body.Read(data)
}

// ReadAll reads the request body to a byte array
func (c *Context) ReadAll() ([]byte, error) {
	return ioutil.ReadAll(c.request.Body)
}

// ReadString reads the request body as a string
func (c *Context) ReadString() (string, error) {
	data, err := c.ReadAll()

	return string(data), err
}

// ReadJson unmarshal the request body to the provided interface.
func (c *Context) ReadJson(data interface{}) error {
	return json.NewDecoder(c.request.Body).Decode(data)
}

// Forms

// ReadForm reads the request body as a form.
func (c *Context) ReadForm() map[string][]string {
	return c.request.Form
}

// ReadFormValue reads the request body as a form and returns the value of the provided key.
func (c *Context) ReadFormValue(key string) string {
	return c.request.FormValue(key)
}

// Headers

// SetHeader sets the provided key and value in the response headers.
func (c *Context) SetHeader(key, value string) {
	c.ResponseWriter().Header().Set(key, value)
}

// SetHeaders overwrites all headers
// This is not recommended for real use
func (c *Context) SetHeaders(headers map[string][]string) {
	c.request.Header = headers
}

// GetHeader gets the provided key from the request headers.
func (c *Context) GetHeader(key string) string {
	return c.Request().Header.Get(key)
}

// GetHeaders gets all the headers from the request.
// 1 header can have an array of values that's why []string is the map's value
func (c *Context) GetHeaders() map[string][]string {
	return c.Request().Header
}

// Params

// GetParam gets the provided key from the request params.
func (c *Context) GetParam(key string) string {
	return c.vars[key]
}

// HasParam checks if the provided key is in the request params.
func (c *Context) HasParam(key string) bool {
	_, ok := c.vars[key]

	return ok
}

// GetParams gets all the params from the request.
func (c *Context) GetParams() map[string]string {
	return c.vars
}

// Query

// GetQuery gets the provided key from the request query.
func (c *Context) GetQuery(key string) string {
	return c.request.URL.Query().Get(key)
}

// Queries gets all the queries from the request.
func (c *Context) Queries() map[string][]string {
	return c.request.URL.Query()
}

// Cookies

// GetCookie gets the provided key from the request cookies.
func (c *Context) GetCookie(key string) (*http.Cookie, error) {
	return c.request.Cookie(key)
}

// SetCookie set a response cookie.
func (c *Context) SetCookie(key string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    value,
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	}

	c.SetResponseCookie(cookie)
}

// SetResponseCookie sets a response cookie.
func (c *Context) SetResponseCookie(cookies ...*http.Cookie) {
	for _, cookie := range cookies {
		http.SetCookie(c.ResponseWriter(), cookie)
	}
}

// Templates

func (c *Context) RenderTemplate(template string, data map[string]drivers.Value) error {
	return c.a.TemplateDriver.RenderTemplate(template, c.responseWriter, data)
}

// For middleware

// Abort stops the current middleware chain.
func (c *Context) Abort() {
	c.abort = true
}

// Aborted returns true if the context was aborted.
func (c *Context) Aborted() bool {
	return c.abort
}
