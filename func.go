package functional

func SliceMap[In any, Out any](in []In, fn func(In) Out) []Out {
	out := make([]Out, 0, len(in))

	for _, v := range in {
		out = append(out, fn(v))
	}

	return out
}
func SliceMapWithError[In any, Out any](in []In, fn func(In) (Out, error)) ([]Out, error) {
	out := make([]Out, 0, len(in))

	for _, v := range in {
		o, err := fn(v)
		if err != nil {
			return nil, err
		}
		out = append(out, o)
	}

	return out, nil
}
func Contain[In comparable](in []In, v In) bool {
	return ContainWithFn(in, func(i In) bool { return i == v })
}

func ContainWithFn[In any](in []In, fn func(In) bool) bool {
	for _, v := range in {
		if fn(v) {
			return true
		}
	}
	return false
}
func SliceFilter[T any](in []T, fn func(T) bool) []T {
	out := make([]T, 0, len(in))

	for _, v := range in {
		if fn(v) {
			out = append(out, v)
		}
	}

	return out
}
func ToLookupTable[KeyType comparable, ElemType any](in []ElemType, keyFn func(ElemType) KeyType) map[KeyType]ElemType {
	m := make(map[KeyType]ElemType, len(in))

	for _, elem := range in {
		key := keyFn(elem)
		m[key] = elem
	}

	return m
}
func Some[T any](in []T, fn func(T) bool) bool {
	return ContainWithFn(in, fn)
}
func All[T any](in []T, fn func(T) bool) bool {
	for _, v := range in {
		if !fn(v) {
			return false
		}
	}
	return true
}
func Reduce[In any, Out any](in []In, fn func(accumulator Out, v In) Out, init Out) Out {
	out := init

	for _, v := range in {
		out = fn(out, v)
	}

	return out
}
func Flat[T any](in [][]T) []T {
	size := 0
	for _, arr := range in {
		size += len(arr)
	}

	out := make([]T, 0, size)
	for _, arr := range in {
		out = append(out, arr...)
	}
	return out
}

func First[In any](in []In, fn func(In) bool) *In {
	for i := range in {
		if fn(in[i]) {
			return &in[i]
		}
	}
	return nil
}

func Last[In any](in []In, fn func(In) bool) *In {
	for i := len(in) - 1; i >= 0; i-- {
		if fn(in[i]) {
			return &in[i]
		}
	}
	return nil
}

func ForEach[T any](in []T, fn func(T) error) error {
	for _, v := range in {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func SliceInsertFirst[T any](slice []T, elem T) []T {
	out := make([]T, 0, len(slice)+1)
	out = append(out, elem)
	out = append(out, slice...)
	return out
}

func SliceInsertLast[T any](slice []T, elem T) []T {
	out := make([]T, 0, len(slice)+1)
	out = append(out, slice...)
	out = append(out, elem)
	return out
}
