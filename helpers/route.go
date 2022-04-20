package helpers

import (
	"github.com/volix-dev/leopard"
)

// Route gets the route's path.
// Args should be key first and value second.
// Having incorrect args will result in a panic.
func Route(name string, args ...string) string {
	route := leopard.Instance.GetRouter().GetRoute(name)

	if route == nil {
		return ""
	}

	url, err := route.URL(args...)

	if err != nil {
		panic(err)
	}

	if url != nil {
		return url.Path
	}

	return ""
}
