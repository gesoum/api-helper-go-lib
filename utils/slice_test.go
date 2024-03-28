package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitToChunk(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	r1 := SplitToChunk(1, []int{1, 2, 3, 4})
	is.Equal(r1, [][]int{{1}, {2}, {3}, {4}})

	r2 := SplitToChunk(2, []string{"a", "b", "c"})
	is.Equal(r2, [][]string{{"a", "b"}, {"c"}})
}
