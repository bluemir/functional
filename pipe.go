package functional

import "fmt"

type PipeFn func(any /* []In */) (any /* []Out */, error)

func MapFn[In, Out any](fn func(In) Out) PipeFn {
	return func(input any /* []In */) (any /* []Out */, error) {
		slice, ok := input.([]In)
		if !ok {
			return nil, fmt.Errorf("MapFn: type assertion failed: expected []%T, got %T", *new(In), input)
		}
		return Map(slice, fn), nil
	}
}

func FilterFn[T any](fn func(T) bool) PipeFn {
	return func(input any /* []T */) (any /* []T */, error) {
		slice, ok := input.([]T)
		if !ok {
			return nil, fmt.Errorf("FilterFn: type assertion failed: expected []%T, got %T", *new(T), input)
		}
		return Filter(slice, fn), nil
	}
}

func MapWithErrorFn[In, Out any](fn func(In) (Out, error)) PipeFn {
	return func(input any /* []In */) (any /* []Out */, error) {
		slice, ok := input.([]In)
		if !ok {
			return nil, fmt.Errorf("MapWithErrorFn: type assertion failed: expected []%T, got %T", *new(In), input)
		}
		return MapWithError(slice, fn)
	}
}


// InsertFirstFn returns a PipeFn that prepends the given element to the slice.
func InsertFirstFn[T any](elem T) PipeFn {
	return func(input any /* []T */) (any /* []T */, error) {
		slice, ok := input.([]T)
		if !ok {
			return nil, fmt.Errorf("InsertFirstFn: type assertion failed: expected []%T, got %T", *new(T), input)
		}
		return InsertFirst(slice, elem), nil
	}
}

// InsertLastFn returns a PipeFn that appends the given element to the slice.
func InsertLastFn[T any](elem T) PipeFn {
	return func(input any /* []T */) (any /* []T */, error) {
		slice, ok := input.([]T)
		if !ok {
			return nil, fmt.Errorf("InsertLastFn: type assertion failed: expected []%T, got %T", *new(T), input)
		}
		return InsertLast(slice, elem), nil
	}
}

// ConsFn is an alias for InsertFirstFn (classic FP "cons" operation).
func ConsFn[T any](elem T) PipeFn {
	return InsertFirstFn(elem)
}

// example
// functional.Pipe[int, string](
//   []int{1, 2, 3},
//   functional.FilterFn(func(i int) bool { return i > 1 }),
//   functional.MapFn(func(i int) string { return strconv.Itoa(i * 10) }),
//   functional.MapWithErrorFn(func(s string) (string, error) { return s + "!", nil }),
// ) // return []string{"20!", "30!"}, nil
func Pipe[In, Out any](input []In, fns ...PipeFn) ([]Out, error) {
	var current any = input
	for _, fn := range fns {
		result, err := fn(current)
		if err != nil {
			return nil, err
		}
		current = result
	}
	result, ok := current.([]Out)
	if !ok {
		return nil, fmt.Errorf("Pipe: type assertion failed: expected []%T, got %T", *new(Out), current)
	}
	return result, nil
}
