package twigDriver

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tyler-sommer/stick"
	"io"
	"leopard/templating/drivers"
	"reflect"
)

type TwigDriver struct {
	env *stick.Env
}

func NewTwigDriver() *TwigDriver {
	return &TwigDriver{
		env: stick.New(nil),
	}
}

func (t *TwigDriver) RenderTemplate(template string, writer io.Writer, data map[string]drivers.Value) error {
	castedData := make(map[string]stick.Value)

	for s, value := range data {
		castedData[s] = stick.Value(value)
	}

	return t.env.Execute(template, writer, castedData)
}

func (t *TwigDriver) Load(path string, router *mux.Router) error {
	t.env.Loader = stick.NewFilesystemLoader(path)

	t.env.Functions["route"] = func(ctx stick.Context, args ...stick.Value) stick.Value {
		if len(args) == 0 {
			return stick.Value("")
		}
		route := router.GetRoute(args[0].(string))

		if route == nil {
			panic("route not found")
		}

		var stringData []string

		for _, v := range args[1:] {

			switch reflect.TypeOf(v).Kind() {
			case reflect.String:
				stringData = append(stringData, v.(string))
				break

			case reflect.Float64:
				stringData = append(stringData, fmt.Sprintf("%g", v.(float64)))
				break
			}
		}

		url, _ := route.URL(stringData...)

		if url != nil {
			return url.Path
		}

		return ""
	}

	t.env.Functions["format"] = func(ctx stick.Context, args ...stick.Value) stick.Value {
		if len(args) < 2 {
			return stick.Value("")
		}

		var data []interface{}

		for _, d := range args {
			data = append(data, d)
		}

		return fmt.Sprintf(args[0].(string), data[1:]...)
	}

	t.env.Tests["route"] = func(ctx stick.Context, val stick.Value, args ...stick.Value) bool {
		return router.GetRoute(stick.CoerceString(val)) != nil
	}

	return nil
}
