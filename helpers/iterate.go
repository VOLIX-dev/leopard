package helpers

type Iterator[V any] interface {
	Next() bool
	Value() V
	ResetIterator()
}

func Iterate[V any](it Iterator[V], f func(value V)) {
	for it.Next() {
		f(it.Value())
	}
}
