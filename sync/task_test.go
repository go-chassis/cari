package sync

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {
	t.Run("resource is string", func(t *testing.T) {
		task, err := NewTask("", "", "", "", "hello")
		if assert.Nil(t, err) {
			assert.Equal(t, []byte("hello"), task.Resource)
		}
	})

	t.Run("resource is []byte", func(t *testing.T) {
		task, err := NewTask("", "", "", "", []byte("hello"))
		if assert.Nil(t, err) {
			assert.Equal(t, []byte("hello"), task.Resource)
		}
	})

	t.Run("resource is interface", func(t *testing.T) {
		r := map[string]string{
			"a": "b",
		}
		result, _ := json.Marshal(r)
		task, err := NewTask("", "", "", "", r)
		if assert.Nil(t, err) {
			assert.Equal(t, result, task.Resource)
		}
	})
}
