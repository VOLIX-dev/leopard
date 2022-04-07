package leopard

import (
	"errors"
	"leopard/templating/drivers"
	"testing"
)

func TestLog(t *testing.T) {
	a, err := New()

	if err != nil {
		panic(err)
	}
	a.GET("/kanker/{name}", "test", func(c *Context) {
		err := c.RenderTemplate("test.twig", map[string]drivers.Value{
			"test": "a",
		})
		if err != nil {
			panic(err)
		}
	})

	a.GET("/error", "test2", func(c *Context) {
		c.Error(errors.New("AAAAAAAAA"))
	})

	a.GET("/panic", "test2", func(c *Context) {
		panic("AAAAAAAAA")
	})

	a.Serve()
}
