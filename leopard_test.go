package leopard

import (
	"fmt"
	"net/http"
	"testing"
)

func TestLog(t *testing.T) {
	a, err := New()

	if err != nil {
		panic(err)
	}
	a.GET("/kanker", func(c *http.Request) {
		fmt.Println("hi")
	})

	a.Serve()
}
