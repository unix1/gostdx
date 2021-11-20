package cons

type cons[T any] func() (T, cons[T])

func List[T any](vals ...T) cons[T] {
	var c cons[T]
	for i := len(vals) - 1; i >= 0; i-- {
		c = Cons(vals[i], c)
	}
	return c
}

func Cons[T any](val T, c cons[T]) cons[T] {
	return func() (T, cons[T]) { return val, c }
}

func Car[T any](c cons[T]) T {
	car, _ := c()
	return car
}

func Cdr[T any](c cons[T]) cons[T] {
	_, cdr := c()
	return cdr
}
