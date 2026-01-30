package functional

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollection_Map(t *testing.T) {
	result, err := From[int, int]([]int{1, 2, 3}).
		Map(func(v any) any { return v.(int) * 2 }).
		ToSlice()

	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestCollection_Filter(t *testing.T) {
	result, err := From[int, int]([]int{1, 2, 3, 4, 5}).
		Filter(func(v any) bool { return v.(int)%2 == 0 }).
		ToSlice()

	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4}, result)
}

func TestCollection_MapWithError(t *testing.T) {
	result, err := From[int, int]([]int{1, 2, 3}).
		MapWithError(func(v any) (any, error) { return v.(int) * 2, nil }).
		ToSlice()

	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestCollection_MapWithError_Error(t *testing.T) {
	result, err := From[int, int]([]int{1, 2, 3}).
		MapWithError(func(v any) (any, error) {
			if v.(int) == 2 {
				return nil, errors.New("error")
			}
			return v.(int) * 2, nil
		}).
		ToSlice()

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCollection_Chaining(t *testing.T) {
	result, err := From[int, int]([]int{1, 2, 3, 4, 5}).
		Filter(func(v any) bool { return v.(int)%2 == 0 }).
		Map(func(v any) any { return v.(int) * 10 }).
		ToSlice()

	assert.NoError(t, err)
	assert.Equal(t, []int{20, 40}, result)
}

func TestCollection_Reduce(t *testing.T) {
	result, err := From[int, int]([]int{1, 2, 3, 4}).
		Reduce(func(acc int, v any) int { return acc + v.(int) }, 0)

	assert.NoError(t, err)
	assert.Equal(t, 10, result)
}

func TestCollection_ForEach(t *testing.T) {
	sum := 0
	err := From[int, int]([]int{1, 2, 3}).
		ForEach(func(v int) error {
			sum += v
			return nil
		})

	assert.NoError(t, err)
	assert.Equal(t, 6, sum)
}

func TestCollection_ForEach_Error(t *testing.T) {
	err := From[int, int]([]int{1, 2, 3}).
		ForEach(func(v int) error {
			if v == 2 {
				return errors.New("error")
			}
			return nil
		})

	assert.Error(t, err)
}

func TestCollection_First(t *testing.T) {
	result, err := From[int, int]([]int{1, 2, 3}).First()

	assert.NoError(t, err)
	assert.Equal(t, 1, *result)
}

func TestCollection_First_Empty(t *testing.T) {
	result, err := From[int, int]([]int{}).First()

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCollection_Last(t *testing.T) {
	result, err := From[int, int]([]int{1, 2, 3}).Last()

	assert.NoError(t, err)
	assert.Equal(t, 3, *result)
}

func TestCollection_Last_Empty(t *testing.T) {
	result, err := From[int, int]([]int{}).Last()

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCollection_Pick(t *testing.T) {
	result, err := From[int, int]([]int{1, 2, 3}).Pick(1)

	assert.NoError(t, err)
	assert.Equal(t, 2, *result)
}

func TestCollection_Pick_Empty(t *testing.T) {
	result, err := From[int, int]([]int{}).Pick(0)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCollection_ErrorPropagation(t *testing.T) {
	result, err := From[int, int]([]int{1, 2, 3}).
		MapWithError(func(v any) (any, error) {
			return nil, errors.New("early error")
		}).
		Map(func(v any) any { return v.(int) * 2 }).
		Filter(func(v any) bool { return true }).
		ToSlice()

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSum(t *testing.T) {
	result, err := From[int, int]([]int{1, 2, 3, 4}).
		Reduce(Sum[int], 0)

	assert.NoError(t, err)
	assert.Equal(t, 10, result)
}
