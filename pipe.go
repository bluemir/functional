package functional

import "fmt"

type PipeFn func(any /* []In */) (any /* []Out */, error)

func Map[In, Out any](fn func(In) Out) PipeFn {
	return func(input any /* []In */) (any /* []Out */, error) {
		slice, ok := input.([]In)
		if !ok {
			return nil, fmt.Errorf("Map: type assertion failed: expected []%T, got %T", *new(In), input)
		}
		return SliceMap(slice, fn), nil
	}
}

func Filter[T any](fn func(T) bool) PipeFn {
	return func(input any /* []T */) (any /* []T */, error) {
		slice, ok := input.([]T)
		if !ok {
			return nil, fmt.Errorf("Filter: type assertion failed: expected []%T, got %T", *new(T), input)
		}
		return SliceFilter(slice, fn), nil
	}
}

func MapWithError[In, Out any](fn func(In) (Out, error)) PipeFn {
	return func(input any /* []In */) (any /* []Out */, error) {
		slice, ok := input.([]In)
		if !ok {
			return nil, fmt.Errorf("MapWithError: type assertion failed: expected []%T, got %T", *new(In), input)
		}
		return SliceMapWithError(slice, fn)
	}
}

// InsertFirst returns a PipeFn that prepends the given element to the slice.
func InsertFirst[T any](elem T) PipeFn {
	return func(input any /* []T */) (any /* []T */, error) {
		slice, ok := input.([]T)
		if !ok {
			return nil, fmt.Errorf("InsertFirst: type assertion failed: expected []%T, got %T", *new(T), input)
		}
		return SliceInsertFirst(slice, elem), nil
	}
}

// InsertLast returns a PipeFn that appends the given element to the slice.
func InsertLast[T any](elem T) PipeFn {
	return func(input any /* []T */) (any /* []T */, error) {
		slice, ok := input.([]T)
		if !ok {
			return nil, fmt.Errorf("InsertLast: type assertion failed: expected []%T, got %T", *new(T), input)
		}
		return SliceInsertLast(slice, elem), nil
	}
}

// Cons is an alias for InsertFirst (classic FP "cons" operation).
func Cons[T any](elem T) PipeFn {
	return InsertFirst[T](elem)
}

// example
// functional.Pipe[int, string](
//   []int{1, 2, 3},
//   functional.Filter(func(i int) bool { return i > 1 }),
//   functional.Map(func(i int) string { return strconv.Itoa(i * 10) }),
//   functional.MapWithError(func(s string) (string, error) { return s + "!", nil }),
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
