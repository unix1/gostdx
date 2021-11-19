package cons

type registers[T any] struct {
	car T
	cdr cons[T]
}

type cons[T any] func() registers[T]

func List[T any](vals ...T) cons[T] {
	var c cons[T]
	for i := len(vals) - 1; i >= 0; i-- {
		c = Cons(vals[i], c)
	}
	return c
}

func Cons[T any](val T, c cons[T]) cons[T] {
	return func() registers[T] {
		return registers[T]{
			car: val,
			cdr: c,
		}
	}
}

func Car[T any](c cons[T]) T {
	return c().car
}

func Cdr[T any](c cons[T]) cons[T] {
	return c().cdr
}
