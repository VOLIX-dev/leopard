package database

import (
	"errors"
	. "github.com/volix-dev/leopard/helpers"
	"reflect"
)

type Model struct {
	fieldMap map[string]field
	indexMap map[string]index
	mapped   bool
}

type field struct {
	Name          string
	Type          reflect.Kind
	PrimaryKey    bool
	Key           bool
	AutoIncrement bool
	NotNull       bool
	Default       string
	Unique        bool
}

type index struct {
	Name   string
	Unique bool
	Fields []string
}

func (m *Model) mapModel() {
	if m.mapped {
		return
	}

	m.mapped = true

	t := reflect.TypeOf(m)

	m.fieldMap = make(map[string]field)
	m.indexMap = make(map[string]index)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if f.Anonymous {
			continue
		}

		fied := field{
			Name: getFieldName(f),
		}

		f.Type.Kind()
	}
}

func getFieldName(f reflect.StructField) string {
	if f.Anonymous {
		return f.Name
	}

	column := f.Tag.Get("column")
	return Ternary(column != "", column, f.Name)
}

func getFieldType(f reflect.StructField) reflect.Kind {
	if f.Anonymous {
		return f.Type.Kind()
	}

	column := f.Tag.Get("type")
	typ, found := kindNames[column]

	if !found {
		panic(errors.New("Invalid type: " + column))
	}

	return Ternary(column != "", typ, f.Type.Kind())
}

func (m *Model) GetPrimaryKey() string {

}

var kindNames = map[string]reflect.Kind{
	"invalid":        reflect.Invalid,
	"bool":           reflect.Bool,
	"int":            reflect.Int,
	"int8":           reflect.Int8,
	"int16":          reflect.Int16,
	"int32":          reflect.Int32,
	"int64":          reflect.Int64,
	"uint":           reflect.Uint,
	"uint8":          reflect.Uint8,
	"uint16":         reflect.Uint16,
	"uint32":         reflect.Uint32,
	"uint64":         reflect.Uint64,
	"uintptr":        reflect.Uintptr,
	"float32":        reflect.Float32,
	"float64":        reflect.Float64,
	"complex64":      reflect.Complex64,
	"complex128":     reflect.Complex128,
	"array":          reflect.Array,
	"chan":           reflect.Chan,
	"func":           reflect.Func,
	"interface":      reflect.Interface,
	"map":            reflect.Map,
	"ptr":            reflect.Ptr,
	"slice":          reflect.Slice,
	"string":         reflect.String,
	"struct":         reflect.Struct,
	"unsafe.Pointer": reflect.UnsafePointer,
}
