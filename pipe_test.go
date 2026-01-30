package functional_test

import (
	"testing"

	"github.com/bluemir/functional"
	"github.com/stretchr/testify/assert"
)

func TestPipe(t *testing.T) {
	result := functional.Pipe[int, int](
		[]int{1, 2, 3},
		functional.MapFn(func(i int) int { return i * 2 }),
	)

	assert.Equal(t, []int{2, 4, 6}, result)
}
