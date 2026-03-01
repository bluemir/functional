package functional

import (
	"context"
	"sync"
)

// elemResult holds the result of processing a single element with its original index.
type elemResult struct {
	index  int
	value  any
	keep   bool
	err    error
}

// executeParallel runs a fused ElemFn segment using a worker pool.
func (seg *segment) executeParallel(items []any, cfg *lazyConfig) ([]any, error) {
	// Check context before starting any work
	select {
	case <-cfg.ctx.Done():
		return nil, cfg.ctx.Err()
	default:
	}

	n := len(items)

	ctx, cancel := context.WithCancel(cfg.ctx)
	defer cancel()

	// Channel for dispatching work items
	jobs := make(chan int, n)
	results := make([]elemResult, n)

	var wg sync.WaitGroup
	wg.Add(cfg.workers)

	// Start workers
	for w := 0; w < cfg.workers; w++ {
		go func() {
			defer wg.Done()
			for idx := range jobs {
				// Check context before processing
				select {
				case <-ctx.Done():
					return
				default:
				}

				current := items[idx]
				keep := true
				var err error
				for _, fn := range seg.elemFns {
					current, keep, err = fn(current)
					if err != nil {
						results[idx] = elemResult{index: idx, err: err}
						cancel()
						return
					}
					if !keep {
						break
					}
				}
				results[idx] = elemResult{index: idx, value: current, keep: keep}
			}
		}()
	}

	// Dispatch jobs
	for i := 0; i < n; i++ {
		select {
		case jobs <- i:
		case <-ctx.Done():
			goto done
		}
	}
done:
	close(jobs)

	wg.Wait()

	// Check if the user-supplied context was cancelled
	if cfg.ctx.Err() != nil {
		return nil, cfg.ctx.Err()
	}

	// Collect results (also detects worker errors)
	if cfg.ordered {
		return collectOrdered(results)
	}
	return collectUnordered(results)
}

// collectOrdered collects results preserving the original input order.
func collectOrdered(results []elemResult) ([]any, error) {
	out := make([]any, 0, len(results))
	for _, r := range results {
		if r.err != nil {
			return nil, r.err
		}
		if r.keep {
			out = append(out, r.value)
		}
	}
	return out, nil
}

// collectUnordered collects results without preserving order.
// Since workers process in arbitrary order, results may appear in any order.
// We still iterate by index but the effective order may differ from input
// when elements are filtered out.
func collectUnordered(results []elemResult) ([]any, error) {
	out := make([]any, 0, len(results))
	for _, r := range results {
		if r.err != nil {
			return nil, r.err
		}
		if r.keep {
			out = append(out, r.value)
		}
	}
	return out, nil
}
