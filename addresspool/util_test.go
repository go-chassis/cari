package addresspool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_removeDuplicates(t *testing.T) {
	assert.Empty(t, removeDuplicates(nil))
	assert.Equal(t, []string{"v1", "v2"}, removeDuplicates([]string{"v1", "v2"}))
	assert.Equal(t, []string{"v1", "v2", "v3"}, removeDuplicates([]string{"v1", "v2", "v2", "v1", "v1", "v3"}))
}
