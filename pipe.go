package functional

type PipeFn func(any) any

func MapFn[In, Out any](fn func(In) Out) PipeFn {
	return func(input any) any {
		return Map(input.([]In), fn) // 라이브러리 내부에서 type assertion
	}
}

func FilterFn[T any](fn func(T) bool) PipeFn {
	return func(input any) any {
		return Filter(input.([]T), fn)
	}
}

func Pipe[In, Out any](input []In, fns ...PipeFn) []Out {
	var current any = input
	for _, fn := range fns {
		current = fn(current)
	}
	return current.([]Out)
}
