package arrutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGroupBy(t *testing.T) {
	list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	_map := GroupBy(list, func(in int) bool {
		return in%2 == 0
	})
	assert.Equal(t, 5, len(_map[true]))
	assert.Equal(t, 5, len(_map[false]))
	assert.ElementsMatch(t, []int{1, 3, 5, 7, 9}, _map[false])
	assert.ElementsMatch(t, []int{2, 4, 6, 8, 10}, _map[true])
}

func TestFilter(t *testing.T) {
	list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	list = Filter(list, func(in int) bool {
		return in%2 == 0
	})
	assert.Equal(t, 5, len(list))
	assert.ElementsMatch(t, []int{2, 4, 6, 8, 10}, list)
}
