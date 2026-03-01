package functional

// ElemFn is an element-level transformation function that unifies map, filter,
// and map+error into a single signature.
//   - output: the transformed element
//   - keep: false to filter out the element
//   - err: non-nil to abort the pipeline
type ElemFn func(elem any) (output any, keep bool, err error)

// LazyMap returns an ElemFn that transforms each element using fn.
func LazyMap[In, Out any](fn func(In) Out) ElemFn {
	return func(elem any) (any, bool, error) {
		return fn(elem.(In)), true, nil
	}
}

// LazyFilter returns an ElemFn that keeps only elements satisfying the predicate.
func LazyFilter[T any](fn func(T) bool) ElemFn {
	return func(elem any) (any, bool, error) {
		v := elem.(T)
		return v, fn(v), nil
	}
}

// LazyMapWithError returns an ElemFn that transforms each element using fn,
// propagating any error to abort the pipeline.
func LazyMapWithError[In, Out any](fn func(In) (Out, error)) ElemFn {
	return func(elem any) (any, bool, error) {
		out, err := fn(elem.(In))
		if err != nil {
			return nil, false, err
		}
		return out, true, nil
	}
}

// LazyFilterMap returns an ElemFn that transforms and optionally filters elements.
// If fn returns false as the second value, the element is excluded.
func LazyFilterMap[In, Out any](fn func(In) (Out, bool)) ElemFn {
	return func(elem any) (any, bool, error) {
		out, keep := fn(elem.(In))
		return out, keep, nil
	}
}

// LazyTap returns an ElemFn that applies a side-effect function to each element
// without modifying it. Useful for logging, debugging, or metrics collection.
//
// When used with WithWorkers(n), fn may be called from multiple goroutines
// concurrently. The caller is responsible for ensuring fn is goroutine-safe.
// WithOrdered(true) guarantees output order but not side-effect invocation order.
func LazyTap[T any](fn func(T)) ElemFn {
	return func(elem any) (any, bool, error) {
		v := elem.(T)
		fn(v)
		return v, true, nil
	}
}

// LazyTapWithError returns an ElemFn that applies a side-effect function to each
// element without modifying it. If fn returns a non-nil error, the pipeline is
// aborted.
//
// When used with WithWorkers(n), fn may be called from multiple goroutines
// concurrently. The caller is responsible for ensuring fn is goroutine-safe.
func LazyTapWithError[T any](fn func(T) error) ElemFn {
	return func(elem any) (any, bool, error) {
		v := elem.(T)
		if err := fn(v); err != nil {
			return nil, false, err
		}
		return v, true, nil
	}
}
