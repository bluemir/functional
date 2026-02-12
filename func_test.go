package functional

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceMap(t *testing.T) {
	result := SliceMap([]int{1, 2, 3}, func(v int) int { return v * 2 })
	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestSliceMap_Empty(t *testing.T) {
	result := SliceMap([]int{}, func(v int) int { return v * 2 })
	assert.Empty(t, result)
}

func TestSliceMap_TypeConvert(t *testing.T) {
	result := SliceMap([]int{1, 2, 3}, func(v int) string { return string(rune('a' + v - 1)) })
	assert.Equal(t, []string{"a", "b", "c"}, result)
}

func TestSliceMapWithError(t *testing.T) {
	result, err := SliceMapWithError([]int{1, 2, 3}, func(v int) (int, error) { return v * 2, nil })
	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestSliceMapWithError_Error(t *testing.T) {
	result, err := SliceMapWithError([]int{1, 2, 3}, func(v int) (int, error) {
		if v == 2 {
			return 0, errors.New("error at 2")
		}
		return v * 2, nil
	})
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestContain(t *testing.T) {
	assert.True(t, Contain([]int{1, 2, 3}, 2))
	assert.False(t, Contain([]int{1, 2, 3}, 4))
}

func TestContainWithFn(t *testing.T) {
	assert.True(t, ContainWithFn([]int{1, 2, 3}, func(v int) bool { return v > 2 }))
	assert.False(t, ContainWithFn([]int{1, 2, 3}, func(v int) bool { return v > 3 }))
}

func TestSliceFilter(t *testing.T) {
	result := SliceFilter([]int{1, 2, 3, 4, 5}, func(v int) bool { return v%2 == 0 })
	assert.Equal(t, []int{2, 4}, result)
}

func TestSliceFilter_Empty(t *testing.T) {
	result := SliceFilter([]int{1, 2, 3}, func(v int) bool { return v > 10 })
	assert.Empty(t, result)
}

func TestToLookupTable(t *testing.T) {
	type item struct {
		ID   int
		Name string
	}
	items := []item{{1, "a"}, {2, "b"}, {3, "c"}}
	result := ToLookupTable(items, func(v item) int { return v.ID })

	assert.Equal(t, item{1, "a"}, result[1])
	assert.Equal(t, item{2, "b"}, result[2])
	assert.Equal(t, item{3, "c"}, result[3])
}

func TestSome(t *testing.T) {
	assert.True(t, Some([]int{1, 2, 3}, func(v int) bool { return v == 2 }))
	assert.False(t, Some([]int{1, 2, 3}, func(v int) bool { return v == 4 }))
}

func TestAll(t *testing.T) {
	assert.True(t, All([]int{2, 4, 6}, func(v int) bool { return v%2 == 0 }))
	assert.False(t, All([]int{1, 2, 3}, func(v int) bool { return v%2 == 0 }))
}

func TestReduce(t *testing.T) {
	result := Reduce([]int{1, 2, 3, 4}, func(acc int, v int) int { return acc + v }, 0)
	assert.Equal(t, 10, result)
}

func TestReduce_WithInit(t *testing.T) {
	result := Reduce([]int{1, 2, 3}, func(acc int, v int) int { return acc + v }, 10)
	assert.Equal(t, 16, result)
}

func TestFlat(t *testing.T) {
	result := Flat([][]int{{1, 2}, {3, 4}, {5}})
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
}

func TestFlat_Empty(t *testing.T) {
	result := Flat([][]int{})
	assert.Empty(t, result)
}

func TestFirst(t *testing.T) {
	result := First([]int{1, 2, 3, 4}, func(v int) bool { return v > 2 })
	assert.NotNil(t, result)
	assert.Equal(t, 3, *result)
}

func TestFirst_NotFound(t *testing.T) {
	result := First([]int{1, 2, 3}, func(v int) bool { return v > 10 })
	assert.Nil(t, result)
}

func TestLast(t *testing.T) {
	result := Last([]int{1, 2, 3, 4}, func(v int) bool { return v > 2 })
	assert.NotNil(t, result)
	assert.Equal(t, 4, *result)
}

func TestLast_NotFound(t *testing.T) {
	result := Last([]int{1, 2, 3}, func(v int) bool { return v > 10 })
	assert.Nil(t, result)
}

func TestForEach(t *testing.T) {
	sum := 0
	err := ForEach([]int{1, 2, 3}, func(v int) error {
		sum += v
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, 6, sum)
}

func TestForEach_Error(t *testing.T) {
	err := ForEach([]int{1, 2, 3}, func(v int) error {
		if v == 2 {
			return errors.New("error at 2")
		}
		return nil
	})
	assert.Error(t, err)
}

func TestSliceInsertFirst(t *testing.T) {
	result := SliceInsertFirst([]int{2, 3, 4}, 1)
	assert.Equal(t, []int{1, 2, 3, 4}, result)
}

func TestSliceInsertFirst_Empty(t *testing.T) {
	result := SliceInsertFirst([]int{}, 1)
	assert.Equal(t, []int{1}, result)
}

func TestSliceInsertFirst_String(t *testing.T) {
	result := SliceInsertFirst([]string{"b", "c"}, "a")
	assert.Equal(t, []string{"a", "b", "c"}, result)
}

func TestSliceInsertLast(t *testing.T) {
	result := SliceInsertLast([]int{1, 2, 3}, 4)
	assert.Equal(t, []int{1, 2, 3, 4}, result)
}

func TestSliceInsertLast_Empty(t *testing.T) {
	result := SliceInsertLast([]int{}, 1)
	assert.Equal(t, []int{1}, result)
}

func TestSliceInsertLast_String(t *testing.T) {
	result := SliceInsertLast([]string{"a", "b"}, "c")
	assert.Equal(t, []string{"a", "b", "c"}, result)
}
