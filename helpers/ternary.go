package helpers

func Ternary[A any](comp bool, a A, b A) A {
	if comp {
		return a
	}
	return b
}
