package functional

import "context"

// LazyOption configures the execution behavior of a LazyPipeline.
type LazyOption func(*lazyConfig)

type lazyConfig struct {
	ctx               context.Context
	workers           int  // 0 or 1 = sequential
	chunkSize         int  // barrier chunk size; 0 = auto
	parallelThreshold int  // below this, use sequential execution (default 1024)
	ordered           bool // preserve order in parallel execution (default true)
}

func defaultConfig() *lazyConfig {
	return &lazyConfig{
		ctx:               context.Background(),
		workers:           0,
		chunkSize:         0,
		parallelThreshold: 1024,
		ordered:           true,
	}
}

// WithWorkers sets the number of parallel workers for element-level stages.
// 0 or 1 means sequential execution.
func WithWorkers(n int) LazyOption {
	return func(c *lazyConfig) {
		c.workers = n
	}
}

// WithContext sets the context for the pipeline execution.
// The context is checked between elements and can cancel parallel workers.
func WithContext(ctx context.Context) LazyOption {
	return func(c *lazyConfig) {
		c.ctx = ctx
	}
}

// WithChunkSize sets the barrier chunk size for parallel execution.
// 0 means automatic sizing.
func WithChunkSize(size int) LazyOption {
	return func(c *lazyConfig) {
		c.chunkSize = size
	}
}

// WithParallelThreshold sets the minimum number of elements required
// for parallel execution. Below this threshold, sequential execution is used.
// Default is 1024.
func WithParallelThreshold(n int) LazyOption {
	return func(c *lazyConfig) {
		c.parallelThreshold = n
	}
}

// WithOrdered controls whether parallel execution preserves input order.
// Default is true.
func WithOrdered(ordered bool) LazyOption {
	return func(c *lazyConfig) {
		c.ordered = ordered
	}
}
