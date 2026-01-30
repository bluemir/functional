package functional

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	result := Keys(m)

	sort.Strings(result)
	assert.Equal(t, []string{"a", "b", "c"}, result)
}

func TestKeys_Empty(t *testing.T) {
	m := map[string]int{}
	result := Keys(m)
	assert.Empty(t, result)
}

func TestValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	result := Values(m)

	sort.Ints(result)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestValues_Empty(t *testing.T) {
	m := map[string]int{}
	result := Values(m)
	assert.Empty(t, result)
}

func TestEntries(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	result := Entries(m)

	assert.Len(t, result, 2)

	found := make(map[string]int)
	for _, p := range result {
		found[p.Key] = p.Value
	}
	assert.Equal(t, m, found)
}

func TestEntries_Empty(t *testing.T) {
	m := map[string]int{}
	result := Entries(m)
	assert.Empty(t, result)
}
