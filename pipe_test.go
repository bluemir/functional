package functional

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipe(t *testing.T) {
	result, err := Pipe[int, int](
		[]int{1, 2, 3},
		Map(func(i int) int { return i * 2 }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestPipeWithFilter(t *testing.T) {
	result, err := Pipe[int, int](
		[]int{1, 2, 3, 4, 5},
		Filter(func(i int) bool { return i%2 == 0 }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4}, result)
}

func TestPipeMapAndFilter(t *testing.T) {
	result, err := Pipe[int, int](
		[]int{1, 2, 3, 4, 5},
		Map(func(i int) int { return i * 2 }),
		Filter(func(i int) bool { return i > 5 }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{6, 8, 10}, result)
}

func TestPipeTypeConversion(t *testing.T) {
	result, err := Pipe[int, string](
		[]int{1, 2, 3},
		Map(func(i int) string { return strconv.Itoa(i) }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []string{"1", "2", "3"}, result)
}

func TestPipeEmptySlice(t *testing.T) {
	result, err := Pipe[int, int](
		[]int{},
		Map(func(i int) int { return i * 2 }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{}, result)
}

func TestPipeTypeAssertionError(t *testing.T) {
	_, err := Pipe[int, string](
		[]int{1, 2, 3},
		Map(func(i int) int { return i * 2 }),
	)

	assert.Error(t, err)
}

func TestPipeChainedMaps(t *testing.T) {
	result, err := Pipe[int, int](
		[]int{1, 2, 3},
		Map(func(i int) int { return i + 1 }),
		Map(func(i int) int { return i * 2 }),
		Map(func(i int) int { return i - 1 }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{3, 5, 7}, result)
}

func TestPipeMultipleTypeChanges(t *testing.T) {
	// int -> string -> string
	result, err := Pipe[int, string](
		[]int{1, 2, 3},
		Map(func(i int) string { return strconv.Itoa(i) }),
		Map(func(s string) string { return s + "!" }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []string{"1!", "2!", "3!"}, result)
}

func TestPipeIntStringInt(t *testing.T) {
	// int -> string -> int
	result, err := Pipe[int, int](
		[]int{1, 2, 3},
		Map(func(i int) string { return strconv.Itoa(i * 10) }),
		Map(func(s string) int {
			n, _ := strconv.Atoi(s)
			return n + 1
		}),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{11, 21, 31}, result)
}

func TestMapWithError(t *testing.T) {
	result, err := Pipe[string, int](
		[]string{"1", "2", "3"},
		MapWithError(func(s string) (int, error) {
			return strconv.Atoi(s)
		}),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestMapWithErrorReturnsError(t *testing.T) {
	_, err := Pipe[string, int](
		[]string{"1", "not a number", "3"},
		MapWithError(func(s string) (int, error) {
			return strconv.Atoi(s)
		}),
	)

	assert.Error(t, err)
}

func TestMapWithErrorInPipeline(t *testing.T) {
	_, err := Pipe[int, int](
		[]int{1, 2, -3, 4},
		MapWithError(func(i int) (int, error) {
			if i < 0 {
				return 0, fmt.Errorf("negative number: %d", i)
			}
			return i * 2, nil
		}),
	)

	assert.ErrorContains(t, err, "negative number: -3")
}

func TestInsertFirst(t *testing.T) {
	result, err := Pipe[int, int](
		[]int{2, 3, 4},
		InsertFirst(1),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, result)
}

func TestInsertFirst_Empty(t *testing.T) {
	result, err := Pipe[int, int](
		[]int{},
		InsertFirst(1),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{1}, result)
}

func TestInsertFirst_TypeAssertionError(t *testing.T) {
	_, err := Pipe[int, string](
		[]int{1, 2, 3},
		Map(func(i int) string { return fmt.Sprintf("%d", i) }),
		InsertFirst(0), // int into []string -> error
	)

	assert.Error(t, err)
}

func TestInsertLast(t *testing.T) {
	result, err := Pipe[int, int](
		[]int{1, 2, 3},
		InsertLast(4),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, result)
}

func TestInsertLast_Empty(t *testing.T) {
	result, err := Pipe[int, int](
		[]int{},
		InsertLast(1),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{1}, result)
}

func TestInsertLast_InPipeline(t *testing.T) {
	result, err := Pipe[int, int](
		[]int{1, 2, 3},
		Filter(func(i int) bool { return i > 1 }),
		InsertLast(99),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{2, 3, 99}, result)
}

func TestCons(t *testing.T) {
	result, err := Pipe[int, int](
		[]int{2, 3, 4},
		Cons(1),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, result)
}

func TestInsertFirst_InPipeline(t *testing.T) {
	result, err := Pipe[int, int](
		[]int{3, 4, 5},
		InsertFirst(2),
		InsertFirst(1),
		InsertLast(6),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
}

// --- Tap Tests ---

func TestPipeTap(t *testing.T) {
	var seen []int
	result, err := Pipe[int, int](
		[]int{1, 2, 3},
		Tap(func(i int) { seen = append(seen, i) }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
	assert.Equal(t, []int{1, 2, 3}, seen)
}

func TestPipeTap_InPipeline(t *testing.T) {
	var tapped []int
	result, err := Pipe[int, string](
		[]int{1, 2, 3, 4, 5},
		Filter(func(i int) bool { return i > 2 }),
		Tap(func(i int) { tapped = append(tapped, i) }),
		Map(func(i int) string { return strconv.Itoa(i * 10) }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []string{"30", "40", "50"}, result)
	assert.Equal(t, []int{3, 4, 5}, tapped)
}

func TestPipeTapWithError_Error(t *testing.T) {
	_, err := Pipe[int, int](
		[]int{1, 2, 3},
		TapWithError(func(i int) error {
			if i == 2 {
				return fmt.Errorf("tap error at %d", i)
			}
			return nil
		}),
	)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "tap error at 2")
}

func TestPipeTap_TypeAssertionError(t *testing.T) {
	_, err := Pipe[int, string](
		[]int{1, 2, 3},
		Map(func(i int) string { return strconv.Itoa(i) }),
		Tap(func(i int) {}), // string slice, but Tap expects int
	)

	assert.Error(t, err)
}

// --- Once Tests ---

func TestPipeOnce(t *testing.T) {
	callCount := 0
	result, err := Pipe[int, int](
		[]int{1, 2, 3},
		Once(func() error { callCount++; return nil }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
	assert.Equal(t, 1, callCount)
}

func TestPipeOnce_InPipeline(t *testing.T) {
	var marker int
	result, err := Pipe[int, int](
		[]int{1, 2, 3, 4, 5},
		Filter(func(i int) bool { return i > 2 }),
		Once(func() error { marker = 42; return nil }),
		Map(func(i int) int { return i + marker }),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{45, 46, 47}, result)
}

func TestPipeOnce_Error(t *testing.T) {
	_, err := Pipe[int, int](
		[]int{1, 2, 3},
		Once(func() error { return fmt.Errorf("once error") }),
	)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "once error")
}

func TestPipeOnceWith(t *testing.T) {
	var sum int
	result, err := Pipe[int, int](
		[]int{1, 2, 3},
		OnceWith[int](func(slice []int) error {
			sum = 0
			for _, v := range slice {
				sum += v
			}
			return nil
		}),
	)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
	assert.Equal(t, 6, sum)
}

func TestPipeOnceWith_TypeAssertionError(t *testing.T) {
	_, err := Pipe[int, string](
		[]int{1, 2, 3},
		Map(func(i int) string { return strconv.Itoa(i) }),
		OnceWith[int](func(slice []int) error { return nil }), // []string, not []int
	)

	assert.Error(t, err)
}
