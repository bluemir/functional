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
		functional.Map(func(i int) int { return i * 2 }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestPipeWithFilter(t *testing.T) {
	result, err := functional.Pipe[int, int](
		[]int{1, 2, 3, 4, 5},
		functional.Filter(func(i int) bool { return i%2 == 0 }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4}, result)
}

func TestPipeMapAndFilter(t *testing.T) {
	result, err := functional.Pipe[int, int](
		[]int{1, 2, 3, 4, 5},
		functional.Map(func(i int) int { return i * 2 }),
		functional.Filter(func(i int) bool { return i > 5 }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{6, 8, 10}, result)
}

func TestPipeTypeConversion(t *testing.T) {
	result, err := functional.Pipe[int, string](
		[]int{1, 2, 3},
		functional.Map(func(i int) string { return strconv.Itoa(i) }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []string{"1", "2", "3"}, result)
}

func TestPipeEmptySlice(t *testing.T) {
	result, err := functional.Pipe[int, int](
		[]int{},
		functional.Map(func(i int) int { return i * 2 }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{}, result)
}

func TestPipeTypeAssertionError(t *testing.T) {
	_, err := functional.Pipe[int, string](
		[]int{1, 2, 3},
		functional.Map(func(i int) int { return i * 2 }),
	)

	assert.Error(t, err)
}

func TestPipeChainedMaps(t *testing.T) {
	result, err := functional.Pipe[int, int](
		[]int{1, 2, 3},
		functional.Map(func(i int) int { return i + 1 }),
		functional.Map(func(i int) int { return i * 2 }),
		functional.Map(func(i int) int { return i - 1 }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{3, 5, 7}, result)
}

func TestPipeMultipleTypeChanges(t *testing.T) {
	// int -> string -> string
	result, err := functional.Pipe[int, string](
		[]int{1, 2, 3},
		functional.Map(func(i int) string { return strconv.Itoa(i) }),
		functional.Map(func(s string) string { return s + "!" }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []string{"1!", "2!", "3!"}, result)
}

func TestPipeIntStringInt(t *testing.T) {
	// int -> string -> int
	result, err := functional.Pipe[int, int](
		[]int{1, 2, 3},
		functional.Map(func(i int) string { return strconv.Itoa(i * 10) }),
		functional.Map(func(s string) int {
			n, _ := strconv.Atoi(s)
			return n + 1
		}),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{11, 21, 31}, result)
}

func TestMapWithError(t *testing.T) {
	result, err := functional.Pipe[string, int](
		[]string{"1", "2", "3"},
		functional.MapWithError(func(s string) (int, error) {
			return strconv.Atoi(s)
		}),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestMapWithErrorReturnsError(t *testing.T) {
	_, err := functional.Pipe[string, int](
		[]string{"1", "not a number", "3"},
		functional.MapWithError(func(s string) (int, error) {
			return strconv.Atoi(s)
		}),
	)

	assert.Error(t, err)
}

func TestMapWithErrorInPipeline(t *testing.T) {
	_, err := functional.Pipe[int, int](
		[]int{1, 2, -3, 4},
		functional.MapWithError(func(i int) (int, error) {
			if i < 0 {
				return 0, fmt.Errorf("negative number: %d", i)
			}
			return i * 2, nil
		}),
	)

	assert.ErrorContains(t, err, "negative number: -3")
}

func TestInsertFirst(t *testing.T) {
	result, err := functional.Pipe[int, int](
		[]int{2, 3, 4},
		functional.InsertFirst(1),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, result)
}

func TestInsertFirst_Empty(t *testing.T) {
	result, err := functional.Pipe[int, int](
		[]int{},
		functional.InsertFirst(1),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{1}, result)
}

func TestInsertFirst_TypeAssertionError(t *testing.T) {
	_, err := functional.Pipe[int, string](
		[]int{1, 2, 3},
		functional.Map(func(i int) string { return fmt.Sprintf("%d", i) }),
		functional.InsertFirst(0), // int into []string -> error
	)

	assert.Error(t, err)
}

func TestInsertLast(t *testing.T) {
	result, err := functional.Pipe[int, int](
		[]int{1, 2, 3},
		functional.InsertLast(4),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, result)
}

func TestInsertLast_Empty(t *testing.T) {
	result, err := functional.Pipe[int, int](
		[]int{},
		functional.InsertLast(1),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{1}, result)
}

func TestInsertLast_InPipeline(t *testing.T) {
	result, err := functional.Pipe[int, int](
		[]int{1, 2, 3},
		functional.Filter(func(i int) bool { return i > 1 }),
		functional.InsertLast(99),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{2, 3, 99}, result)
}

func TestCons(t *testing.T) {
	result, err := functional.Pipe[int, int](
		[]int{2, 3, 4},
		functional.Cons(1),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, result)
}

func TestInsertFirst_InPipeline(t *testing.T) {
	result, err := functional.Pipe[int, int](
		[]int{3, 4, 5},
		functional.InsertFirst(2),
		functional.InsertFirst(1),
		functional.InsertLast(6),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
}
