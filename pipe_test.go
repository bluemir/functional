package functional_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/bluemir/functional"
	"github.com/stretchr/testify/assert"
)

func TestPipe(t *testing.T) {
	result, err := functional.Pipe[int, int](
		[]int{1, 2, 3},
		functional.MapFn(func(i int) int { return i * 2 }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestPipeWithFilter(t *testing.T) {
	result, err := functional.Pipe[int, int](
		[]int{1, 2, 3, 4, 5},
		functional.FilterFn(func(i int) bool { return i%2 == 0 }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4}, result)
}

func TestPipeMapAndFilter(t *testing.T) {
	result, err := functional.Pipe[int, int](
		[]int{1, 2, 3, 4, 5},
		functional.MapFn(func(i int) int { return i * 2 }),
		functional.FilterFn(func(i int) bool { return i > 5 }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{6, 8, 10}, result)
}

func TestPipeTypeConversion(t *testing.T) {
	result, err := functional.Pipe[int, string](
		[]int{1, 2, 3},
		functional.MapFn(func(i int) string { return strconv.Itoa(i) }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []string{"1", "2", "3"}, result)
}

func TestPipeEmptySlice(t *testing.T) {
	result, err := functional.Pipe[int, int](
		[]int{},
		functional.MapFn(func(i int) int { return i * 2 }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{}, result)
}

func TestPipeTypeAssertionError(t *testing.T) {
	_, err := functional.Pipe[int, string](
		[]int{1, 2, 3},
		functional.MapFn(func(i int) int { return i * 2 }),
	)

	assert.Error(t, err)
}

func TestPipeChainedMaps(t *testing.T) {
	result, err := functional.Pipe[int, int](
		[]int{1, 2, 3},
		functional.MapFn(func(i int) int { return i + 1 }),
		functional.MapFn(func(i int) int { return i * 2 }),
		functional.MapFn(func(i int) int { return i - 1 }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{3, 5, 7}, result)
}

func TestPipeMultipleTypeChanges(t *testing.T) {
	// int -> string -> string
	result, err := functional.Pipe[int, string](
		[]int{1, 2, 3},
		functional.MapFn(func(i int) string { return strconv.Itoa(i) }),
		functional.MapFn(func(s string) string { return s + "!" }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []string{"1!", "2!", "3!"}, result)
}

func TestPipeIntStringInt(t *testing.T) {
	// int -> string -> int
	result, err := functional.Pipe[int, int](
		[]int{1, 2, 3},
		functional.MapFn(func(i int) string { return strconv.Itoa(i * 10) }),
		functional.MapFn(func(s string) int {
			n, _ := strconv.Atoi(s)
			return n + 1
		}),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{11, 21, 31}, result)
}

func TestMapWithErrorFn(t *testing.T) {
	result, err := functional.Pipe[string, int](
		[]string{"1", "2", "3"},
		functional.MapWithErrorFn(func(s string) (int, error) {
			return strconv.Atoi(s)
		}),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestMapWithErrorFnReturnsError(t *testing.T) {
	_, err := functional.Pipe[string, int](
		[]string{"1", "not a number", "3"},
		functional.MapWithErrorFn(func(s string) (int, error) {
			return strconv.Atoi(s)
		}),
	)

	assert.Error(t, err)
}

func TestMapWithErrorFnInPipeline(t *testing.T) {
	_, err := functional.Pipe[int, int](
		[]int{1, 2, -3, 4},
		functional.MapWithErrorFn(func(i int) (int, error) {
			if i < 0 {
				return 0, fmt.Errorf("negative number: %d", i)
			}
			return i * 2, nil
		}),
	)

	assert.ErrorContains(t, err, "negative number: -3")
}
