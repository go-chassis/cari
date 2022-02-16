/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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

	t.Run("resource is string, domain and project is empty, domain and project should return default", func(t *testing.T) {
		task, err := NewTask("", "", "", "", "hello")
		if assert.Nil(t, err) {
			assert.Equal(t, []byte("hello"), task.Resource)
			assert.Equal(t, Default, task.Domain)
			assert.Equal(t, Default, task.Project)
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
