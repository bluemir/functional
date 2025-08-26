package functional

import (
	"github.com/pkg/errors"
	"golang.org/x/exp/constraints"
)

// example
// functional.From[string, float64](arr).
//   Filter(func(str any) bool{ str.(string) != "" }).
//   MapWithError(func(str any) (any, error) { return strconv.Atoi(str.(string)) /* return int,error */ }).
//   Map(func(i any) any { return i / 2.0 /* return float64 */}).
//   ToSlice() // return []float64, error

// NOTE
// 현재 golang 의 generic 한계상 중간 단계는 any 를 쓸수 밖에 없음.
// type 은 개발자의 책임으로 남겨둠
// (struct나 interface 로 감싸도 되지만 어차피 개발자가 직접 확인해야 해서 직접 타입 추론하는 것보다 간단하거나 안전하지 않음)

type Collection[T any] struct {
	items []any
	err   error
}

func From[From, To any](from []From) Collection[To] {
	return Collection[To]{
		items: Map(from, func(v From) any { return v }),
	}
}

func (col Collection[T]) Map(fn func(any) any) Collection[T] { // if want to convert type. use or 'any'
	out := []any{}

	// TODO concurrent
	for _, v := range col.items {
		out = append(out, fn(v))
	}

	return Collection[T]{
		items: out,
		err:   col.err,
	}
}
func (col Collection[T]) MapWithError(fn func(any) (any, error)) Collection[T] {
	out := []any{}

	for _, v := range col.items {
		r, err := fn(v)
		if err != nil {
			return Collection[T]{
				err: col.err,
			}
		}
		out = append(out, r)
	}

	return Collection[T]{
		items: out,
		err:   col.err,
	}

}

func (col Collection[T]) Filter(fn func(any) bool) Collection[T] {
	out := []any{}

	for _, v := range col.items {
		if fn(v) {
			out = append(out, v)
		}
	}

	return Collection[T]{
		items: out,
		err:   col.err,
	}
}

// aggregator
func (c Collection[T]) ToSlice() ([]T, error) {
	if c.err != nil {
		return nil, c.err
	}
	arr, err := MapWithError(c.items, func(v any) (T, error) {
		a, ok := v.(T)
		if !ok {
			return a, errors.Errorf("type convert failed")
		}
		return a, nil
	})
	if err != nil {
		return nil, err
	}
	return arr, nil
}
func (c Collection[T]) Reduce(fn func(accumulator T, current any) T, initialValue T) (T, error) {
	if c.err != nil {
		return initialValue, c.err
	}

	v := initialValue

	for _, item := range c.items {
		v = fn(v, item)
	}

	return v, c.err
}

// example for reduce
func Sum[T constraints.Integer | constraints.Float](acc T, v any) T {
	return acc + v.(T)
}

func (c Collection[T]) First() (*T, error) {
	if c.err != nil {
		return nil, c.err
	}
	if len(c.items) == 0 {
		return nil, errors.Errorf("collection is empty")
	}

	v := c.items[0].(T)

	return &v, nil
}
func (c Collection[T]) Last() (*T, error) {
	if c.err != nil {
		return nil, c.err
	}
	if len(c.items) == 0 {
		return nil, errors.Errorf("collection is empty")
	}

	v := c.items[len(c.items)-1].(T)

	return &v, nil
}
