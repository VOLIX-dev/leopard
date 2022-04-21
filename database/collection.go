package database

type KeyType interface {
	int | string | int8 | int16 | int32 | int64 | float32 | float64 | complex64 | complex128 | bool | uint | uint8 | uint16 | uint32 | uint64 | uintptr
}

type Collection[Key KeyType, Value any] struct {
	index        int
	counter      int
	data         map[Key]Value
	cachedValues []Value
}

func NewCollection[Key KeyType, Value any](values map[Key]Value) Collection[Key, Value] {
	return Collection[Key, Value]{
		index:   -1,
		data:    values,
		counter: len(values),
	}
}

func (c *Collection[Key, Value]) Next() bool {
	c.index++

	return c.index < len(c.data)
}

func (c *Collection[Key, Value]) Value() Value {
	return c.GetValues()[c.index]
}

func (c *Collection[Key, Value]) ResetIterator() {
	c.index = -1
}

func (c *Collection[Key, Value]) GetValues() []Value {
	if len(c.cachedValues) > 0 {
		return c.cachedValues
	}

	values := make([]Value, len(c.data))

	for _, value := range c.data {
		values = append(values, value)
	}

	return values
}

func (c *Collection[Key, Value]) onUpdate() {
	c.cachedValues = []Value{}
}

func (c *Collection[Key, Value]) Add(key Key, value Value) {
	c.data[key] = value
	c.onUpdate()
	c.counter++
}

func (c *Collection[Key, Value]) AddI(value Value) {
	c.data[c.counter] = value
	c.counter++
	c.onUpdate()
}

func (c *Collection[Key, Value]) Remove(key Key) {
	delete(c.data, key)
	c.onUpdate()
}

func (c *Collection[Key, Value]) RemoveAll() {
	c.data = map[Key]Value{}
	c.onUpdate()
}

func (c *Collection[Key, Value]) Get(key Key) Value {
	return c.data[key]
}
