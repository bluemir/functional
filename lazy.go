package functional

import "fmt"

// BarrierFn wraps a PipeFn with type-safe []any ↔ []T converters,
// capturing type information at construction time via generics
// so that no reflect is needed at execution time.
type BarrierFn struct {
	run func([]any) ([]any, error)
}

// Barrier wraps a PipeFn for use in a LazyPipeline.
// The type parameters capture the input and output element types,
// enabling type-safe conversion between []any and typed slices
// without reflect.
func Barrier[In, Out any](fn PipeFn) BarrierFn {
	return BarrierFn{
		run: func(items []any) ([]any, error) {
			// []any → []In
			typed := make([]In, len(items))
			for i, v := range items {
				typed[i] = v.(In)
			}

			// Run PipeFn
			result, err := fn(typed)
			if err != nil {
				return nil, err
			}

			// []Out → []any
			outSlice := result.([]Out)
			out := make([]any, len(outSlice))
			for i, v := range outSlice {
				out[i] = v
			}
			return out, nil
		},
	}
}

// stage represents a single step in a lazy pipeline.
// Exactly one of elemFn or barrierFn is non-nil.
type stage struct {
	elemFn    ElemFn                      // element-level (fusible)
	barrierFn func([]any) ([]any, error) // slice-level (barrier)
}

// segment is a group of consecutive ElemFn stages that can be fused
// into a single loop, or a single barrier.
type segment struct {
	elemFns   []ElemFn                    // non-empty for fusible segments
	barrierFn func([]any) ([]any, error) // non-nil for barrier segments
}

// LazyPipeline is a deferred execution pipeline that collects stages
// and executes them on Run(). Consecutive element-level stages are fused
// into a single loop to avoid intermediate slice allocations.
type LazyPipeline[In, Out any] struct {
	input  []In
	stages []stage
}

// Lazy creates a new lazy pipeline with the given input slice.
func Lazy[In, Out any](input []In) *LazyPipeline[In, Out] {
	return &LazyPipeline[In, Out]{
		input: input,
	}
}

// Elem appends element-level transformation stages to the pipeline.
// Consecutive Elem stages are fused into a single loop during execution.
func (lp *LazyPipeline[In, Out]) Elem(fns ...ElemFn) *LazyPipeline[In, Out] {
	for _, fn := range fns {
		lp.stages = append(lp.stages, stage{elemFn: fn})
	}
	return lp
}

// Pipe appends slice-level transformation stages (barriers) to the pipeline.
// Each barrier forces materialization of preceding element-level stages.
// Use Barrier[In, Out](pipeFn) to wrap an existing PipeFn.
func (lp *LazyPipeline[In, Out]) Pipe(fns ...BarrierFn) *LazyPipeline[In, Out] {
	for _, fn := range fns {
		lp.stages = append(lp.stages, stage{barrierFn: fn.run})
	}
	return lp
}

// Once appends a barrier that executes fn exactly once (not per-element).
// It forces materialization of preceding Elem stages.
// Useful for side-effects (DB queries, logging) whose results
// are captured via closures for subsequent stages.
func (lp *LazyPipeline[In, Out]) Once(fn func() error) *LazyPipeline[In, Out] {
	lp.stages = append(lp.stages, stage{
		barrierFn: func(items []any) ([]any, error) {
			if err := fn(); err != nil {
				return nil, err
			}
			return items, nil
		},
	})
	return lp
}

// Run executes the lazy pipeline and returns the final result.
func (lp *LazyPipeline[In, Out]) Run(opts ...LazyOption) ([]Out, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	segments := buildSegments(lp.stages)

	// Convert input to []any
	items := make([]any, len(lp.input))
	for i, v := range lp.input {
		items[i] = v
	}

	// Execute each segment
	var err error
	for _, seg := range segments {
		if seg.barrierFn != nil {
			items, err = seg.barrierFn(items)
		} else {
			items, err = seg.execute(items, cfg)
		}
		if err != nil {
			return nil, err
		}
	}

	// Convert []any to []Out
	result := make([]Out, len(items))
	for i, v := range items {
		out, ok := v.(Out)
		if !ok {
			return nil, fmt.Errorf("Lazy: final type assertion failed: expected %T, got %T at index %d", *new(Out), v, i)
		}
		result[i] = out
	}
	return result, nil
}

// buildSegments groups consecutive stages into fusible segments and barriers.
func buildSegments(stages []stage) []segment {
	var segments []segment
	var currentElems []ElemFn

	for _, s := range stages {
		if s.barrierFn != nil {
			// Flush accumulated ElemFns as a segment
			if len(currentElems) > 0 {
				segments = append(segments, segment{elemFns: currentElems})
				currentElems = nil
			}
			segments = append(segments, segment{barrierFn: s.barrierFn})
		} else {
			currentElems = append(currentElems, s.elemFn)
		}
	}

	// Flush remaining ElemFns
	if len(currentElems) > 0 {
		segments = append(segments, segment{elemFns: currentElems})
	}

	return segments
}

// execute runs a fusible segment (element-level) either sequentially or in parallel.
func (seg *segment) execute(items []any, cfg *lazyConfig) ([]any, error) {
	if cfg.workers > 1 && len(items) >= cfg.parallelThreshold {
		return seg.executeParallel(items, cfg)
	}
	return seg.executeSequential(items, cfg)
}

// executeSequential runs fused ElemFn loop: for each element, apply all ElemFns.
func (seg *segment) executeSequential(items []any, cfg *lazyConfig) ([]any, error) {
	result := make([]any, 0, len(items))
	for _, item := range items {
		// Check context cancellation
		select {
		case <-cfg.ctx.Done():
			return nil, cfg.ctx.Err()
		default:
		}

		current := item
		keep := true
		var err error
		for _, fn := range seg.elemFns {
			current, keep, err = fn(current)
			if err != nil {
				return nil, err
			}
			if !keep {
				break
			}
		}
		if keep {
			result = append(result, current)
		}
	}
	return result, nil
}
