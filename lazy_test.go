package functional

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// --- Sequential / Fusion Tests ---

func TestLazyBasicMap(t *testing.T) {
	result, err := Lazy[int, string]([]int{1, 2, 3}).
		Elem(LazyMap[int, string](func(i int) string { return strconv.Itoa(i * 10) })).
		Run()

	assert.NoError(t, err)
	assert.Equal(t, []string{"10", "20", "30"}, result)
}

func TestLazyBasicFilter(t *testing.T) {
	result, err := Lazy[int, int]([]int{1, 2, 3, 4, 5}).
		Elem(LazyFilter[int](func(i int) bool { return i > 3 })).
		Run()

	assert.NoError(t, err)
	assert.Equal(t, []int{4, 5}, result)
}

func TestLazyFusedMapAndFilter(t *testing.T) {
	// Filter then map: fused into a single loop
	result, err := Lazy[int, string]([]int{1, 2, 3, 4, 5}).
		Elem(
			LazyFilter[int](func(i int) bool { return i > 2 }),
			LazyMap[int, string](func(i int) string { return strconv.Itoa(i) }),
		).
		Run()

	assert.NoError(t, err)
	assert.Equal(t, []string{"3", "4", "5"}, result)
}

func TestLazyMapWithError_Success(t *testing.T) {
	result, err := Lazy[int, string]([]int{1, 2, 3}).
		Elem(LazyMapWithError[int, string](func(i int) (string, error) {
			return strconv.Itoa(i), nil
		})).
		Run()

	assert.NoError(t, err)
	assert.Equal(t, []string{"1", "2", "3"}, result)
}

func TestLazyMapWithError_Error(t *testing.T) {
	_, err := Lazy[int, string]([]int{1, 2, 3}).
		Elem(LazyMapWithError[int, string](func(i int) (string, error) {
			if i == 2 {
				return "", fmt.Errorf("error at %d", i)
			}
			return strconv.Itoa(i), nil
		})).
		Run()

	assert.Error(t, err)
	assert.ErrorContains(t, err, "error at 2")
}

func TestLazyFilterMap(t *testing.T) {
	result, err := Lazy[int, string]([]int{1, 2, 3, 4, 5}).
		Elem(LazyFilterMap[int, string](func(i int) (string, bool) {
			if i%2 == 0 {
				return strconv.Itoa(i * 10), true
			}
			return "", false
		})).
		Run()

	assert.NoError(t, err)
	assert.Equal(t, []string{"20", "40"}, result)
}

func TestLazyEmptyInput(t *testing.T) {
	result, err := Lazy[int, int]([]int{}).
		Elem(LazyMap[int, int](func(i int) int { return i * 2 })).
		Run()

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestLazyNoStages(t *testing.T) {
	result, err := Lazy[int, int]([]int{1, 2, 3}).
		Run()

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestLazyTypeConversion(t *testing.T) {
	// int -> string -> int
	result, err := Lazy[int, int]([]int{1, 2, 3}).
		Elem(
			LazyMap[int, string](func(i int) string { return strconv.Itoa(i * 10) }),
			LazyMap[string, int](func(s string) int {
				v, _ := strconv.Atoi(s)
				return v + 1
			}),
		).
		Run()

	assert.NoError(t, err)
	assert.Equal(t, []int{11, 21, 31}, result)
}

func TestLazyChainedMaps(t *testing.T) {
	result, err := Lazy[int, int]([]int{1, 2, 3}).
		Elem(
			LazyMap[int, int](func(i int) int { return i * 2 }),
			LazyMap[int, int](func(i int) int { return i + 1 }),
			LazyMap[int, int](func(i int) int { return i * 3 }),
		).
		Run()

	assert.NoError(t, err)
	// (1*2+1)*3=9, (2*2+1)*3=15, (3*2+1)*3=21
	assert.Equal(t, []int{9, 15, 21}, result)
}

// --- Barrier (PipeFn) Tests ---

func TestLazyWithBarrier(t *testing.T) {
	result, err := Lazy[int, int]([]int{1, 2, 3}).
		Elem(LazyMap[int, int](func(i int) int { return i * 2 })).
		Pipe(Barrier[int, int](InsertFirst(0))).
		Run()

	assert.NoError(t, err)
	assert.Equal(t, []int{0, 2, 4, 6}, result)
}

func TestLazyMixedElemAndPipe(t *testing.T) {
	// Elem -> Pipe barrier -> Elem
	result, err := Lazy[int, int]([]int{1, 2, 3}).
		Elem(LazyMap[int, int](func(i int) int { return i * 2 })).
		Pipe(Barrier[int, int](InsertFirst(0))).
		Elem(LazyFilter[int](func(i int) bool { return i > 3 })).
		Run()

	assert.NoError(t, err)
	// After map: [2, 4, 6], after InsertFirst(0): [0, 2, 4, 6], after filter >3: [4, 6]
	assert.Equal(t, []int{4, 6}, result)
}

func TestLazyPipeOnly(t *testing.T) {
	result, err := Lazy[int, int]([]int{1, 2, 3}).
		Pipe(Barrier[int, int](InsertFirst(0))).
		Pipe(Barrier[int, int](InsertLast(4))).
		Run()

	assert.NoError(t, err)
	assert.Equal(t, []int{0, 1, 2, 3, 4}, result)
}

func TestLazyBarrierWithEmptySlice(t *testing.T) {
	result, err := Lazy[int, int]([]int{}).
		Elem(LazyFilter[int](func(i int) bool { return i > 0 })).
		Pipe(Barrier[int, int](InsertFirst(99))).
		Run()

	assert.NoError(t, err)
	assert.Equal(t, []int{99}, result)
}

func TestLazyMultipleBarriers(t *testing.T) {
	result, err := Lazy[int, int]([]int{3, 1, 2}).
		Pipe(Barrier[int, int](InsertFirst(0))).
		Elem(LazyMap[int, int](func(i int) int { return i + 10 })).
		Pipe(Barrier[int, int](InsertLast(99))).
		Run()

	assert.NoError(t, err)
	// After InsertFirst(0): [0,3,1,2], after map +10: [10,13,11,12], after InsertLast(99): [10,13,11,12,99]
	assert.Equal(t, []int{10, 13, 11, 12, 99}, result)
}

// --- Pipe result consistency tests ---

func TestLazyMatchesPipeResult(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}

	// Eager pipe
	eager, err := Pipe[int, string](input,
		Filter[int](func(i int) bool { return i > 1 }),
		Map[int, string](func(i int) string { return strconv.Itoa(i * 10) }),
	)
	assert.NoError(t, err)

	// Lazy pipe
	lazy, err := Lazy[int, string](input).
		Elem(
			LazyFilter[int](func(i int) bool { return i > 1 }),
			LazyMap[int, string](func(i int) string { return strconv.Itoa(i * 10) }),
		).
		Run()
	assert.NoError(t, err)

	assert.Equal(t, eager, lazy)
}

// --- Parallel Tests ---

func TestLazyParallelBasicMap(t *testing.T) {
	input := make([]int, 100)
	for i := range input {
		input[i] = i
	}

	result, err := Lazy[int, int](input).
		Elem(LazyMap[int, int](func(i int) int { return i * 2 })).
		Run(WithWorkers(4), WithParallelThreshold(10))

	assert.NoError(t, err)
	assert.Len(t, result, 100)
	for i, v := range result {
		assert.Equal(t, i*2, v)
	}
}

func TestLazyParallelOrdered(t *testing.T) {
	input := make([]int, 200)
	for i := range input {
		input[i] = i
	}

	result, err := Lazy[int, int](input).
		Elem(LazyMap[int, int](func(i int) int { return i + 1 })).
		Run(WithWorkers(8), WithParallelThreshold(10), WithOrdered(true))

	assert.NoError(t, err)
	expected := make([]int, 200)
	for i := range expected {
		expected[i] = i + 1
	}
	assert.Equal(t, expected, result)
}

func TestLazyParallelWithFilter(t *testing.T) {
	input := make([]int, 100)
	for i := range input {
		input[i] = i
	}

	result, err := Lazy[int, int](input).
		Elem(
			LazyFilter[int](func(i int) bool { return i%2 == 0 }),
			LazyMap[int, int](func(i int) int { return i * 3 }),
		).
		Run(WithWorkers(4), WithParallelThreshold(10), WithOrdered(true))

	assert.NoError(t, err)
	assert.Len(t, result, 50)
	for idx, v := range result {
		assert.Equal(t, idx*2*3, v)
	}
}

func TestLazyParallelErrorPropagation(t *testing.T) {
	input := make([]int, 100)
	for i := range input {
		input[i] = i
	}

	_, err := Lazy[int, int](input).
		Elem(LazyMapWithError[int, int](func(i int) (int, error) {
			if i == 50 {
				return 0, fmt.Errorf("error at 50")
			}
			return i, nil
		})).
		Run(WithWorkers(4), WithParallelThreshold(10))

	assert.Error(t, err)
	assert.ErrorContains(t, err, "error at 50")
}

func TestLazyParallelContextCancellation(t *testing.T) {
	input := make([]int, 1000)
	for i := range input {
		input[i] = i
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := Lazy[int, int](input).
		Elem(LazyMap[int, int](func(i int) int {
			time.Sleep(time.Millisecond) // Simulate work
			return i
		})).
		Run(WithWorkers(4), WithParallelThreshold(10), WithContext(ctx))

	assert.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestLazyParallelBelowThreshold(t *testing.T) {
	// With 10 items and threshold 100, should run sequentially
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	result, err := Lazy[int, int](input).
		Elem(LazyMap[int, int](func(i int) int { return i * 2 })).
		Run(WithWorkers(4), WithParallelThreshold(100))

	assert.NoError(t, err)
	expected := []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}
	assert.Equal(t, expected, result)
}

func TestLazySequentialContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := Lazy[int, int]([]int{1, 2, 3}).
		Elem(LazyMap[int, int](func(i int) int { return i })).
		Run(WithContext(ctx))

	assert.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestLazyParallelWithBarrier(t *testing.T) {
	input := make([]int, 50)
	for i := range input {
		input[i] = i
	}

	result, err := Lazy[int, int](input).
		Elem(LazyMap[int, int](func(i int) int { return i * 2 })).
		Pipe(Barrier[int, int](InsertFirst(-1))).
		Elem(LazyFilter[int](func(i int) bool { return i >= 0 })).
		Run(WithWorkers(4), WithParallelThreshold(10))

	assert.NoError(t, err)
	// After map: [0,2,4,...,98], InsertFirst(-1): [-1,0,2,...,98]
	// After filter >=0: [0,2,4,...,98]
	assert.Len(t, result, 50)
	assert.Equal(t, 0, result[0])
}

func TestLazyFilterAll(t *testing.T) {
	result, err := Lazy[int, int]([]int{1, 2, 3}).
		Elem(LazyFilter[int](func(i int) bool { return false })).
		Run()

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestLazyFilterNone(t *testing.T) {
	result, err := Lazy[int, int]([]int{1, 2, 3}).
		Elem(LazyFilter[int](func(i int) bool { return true })).
		Run()

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
}

// --- Tap Tests ---

func TestLazyTap(t *testing.T) {
	var seen []int
	result, err := Lazy[int, int]([]int{1, 2, 3}).
		Elem(LazyTap[int](func(i int) { seen = append(seen, i) })).
		Run()

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
	assert.Equal(t, []int{1, 2, 3}, seen)
}

func TestLazyTap_InFusedChain(t *testing.T) {
	var tapped []int
	result, err := Lazy[int, string]([]int{1, 2, 3, 4, 5}).
		Elem(
			LazyFilter[int](func(i int) bool { return i > 2 }),
			LazyTap[int](func(i int) { tapped = append(tapped, i) }),
			LazyMap[int, string](func(i int) string { return strconv.Itoa(i * 10) }),
		).
		Run()

	assert.NoError(t, err)
	assert.Equal(t, []string{"30", "40", "50"}, result)
	// Tap should only see elements that passed the filter (3, 4, 5)
	assert.Equal(t, []int{3, 4, 5}, tapped)
}

func TestLazyTap_EmptyInput(t *testing.T) {
	called := false
	result, err := Lazy[int, int]([]int{}).
		Elem(LazyTap[int](func(i int) { called = true })).
		Run()

	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.False(t, called)
}

func TestLazyTapWithError_Success(t *testing.T) {
	var seen []int
	result, err := Lazy[int, int]([]int{1, 2, 3}).
		Elem(LazyTapWithError[int](func(i int) error {
			seen = append(seen, i)
			return nil
		})).
		Run()

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
	assert.Equal(t, []int{1, 2, 3}, seen)
}

func TestLazyTapWithError_Error(t *testing.T) {
	_, err := Lazy[int, int]([]int{1, 2, 3}).
		Elem(LazyTapWithError[int](func(i int) error {
			if i == 2 {
				return fmt.Errorf("tap error at %d", i)
			}
			return nil
		})).
		Run()

	assert.Error(t, err)
	assert.ErrorContains(t, err, "tap error at 2")
}

func TestLazyTap_Parallel(t *testing.T) {
	input := make([]int, 100)
	for i := range input {
		input[i] = i
	}

	var mu sync.Mutex
	count := 0
	result, err := Lazy[int, int](input).
		Elem(LazyTap[int](func(i int) {
			mu.Lock()
			count++
			mu.Unlock()
		})).
		Run(WithWorkers(4), WithParallelThreshold(10), WithOrdered(true))

	assert.NoError(t, err)
	assert.Len(t, result, 100)
	assert.Equal(t, 100, count)
	// Verify output is unchanged and ordered
	for i, v := range result {
		assert.Equal(t, i, v)
	}
}

// --- Once Tests ---

func TestLazyOnce(t *testing.T) {
	callCount := 0
	result, err := Lazy[int, int]([]int{1, 2, 3}).
		Once(func() error { callCount++; return nil }).
		Run()

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
	assert.Equal(t, 1, callCount)
}

func TestLazyOnce_BetweenElems(t *testing.T) {
	divisor := 0
	result, err := Lazy[int, int]([]int{10, 20, 30}).
		Elem(LazyMap[int, int](func(i int) int { return i * 2 })).
		Once(func() error { divisor = 4; return nil }).
		Elem(LazyMap[int, int](func(i int) int { return i / divisor })).
		Run()

	assert.NoError(t, err)
	// After map *2: [20, 40, 60], divisor set to 4, after /4: [5, 10, 15]
	assert.Equal(t, []int{5, 10, 15}, result)
}

func TestLazyOnce_Error(t *testing.T) {
	_, err := Lazy[int, int]([]int{1, 2, 3}).
		Once(func() error { return fmt.Errorf("once error") }).
		Run()

	assert.Error(t, err)
	assert.ErrorContains(t, err, "once error")
}

func TestLazyOnce_EmptyInput(t *testing.T) {
	callCount := 0
	result, err := Lazy[int, int]([]int{}).
		Once(func() error { callCount++; return nil }).
		Run()

	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, 1, callCount)
}

func TestLazyOnce_WithOnceWith(t *testing.T) {
	var sum int
	result, err := Lazy[int, int]([]int{1, 2, 3}).
		Elem(LazyMap[int, int](func(i int) int { return i * 10 })).
		Pipe(Barrier[int, int](OnceWith[int](func(slice []int) error {
			sum = 0
			for _, v := range slice {
				sum += v
			}
			return nil
		}))).
		Run()

	assert.NoError(t, err)
	assert.Equal(t, []int{10, 20, 30}, result)
	assert.Equal(t, 60, sum)
}
