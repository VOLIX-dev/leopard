package leopard

import (
	"errors"
	"testing"
)

func TestLog(t *testing.T) {
	a, err := New()

	if err != nil {
		panic(err)
	}
	a.GET("/kanker/{name}", func(c *Context) {
		c.WriteStringF("Hello %s", c.GetParam("name"))
	})

	a.GET("/error", func(c *Context) {
		c.Error(errors.New("AAAAAAAAA"))
	})

	a.GET("/panic", func(c *Context) {
		panic("AAAAAAAAA")
	})

	a.Serve()
}
