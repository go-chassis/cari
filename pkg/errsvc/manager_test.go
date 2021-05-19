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
package errsvc_test

import (
	"github.com/go-chassis/cari/pkg/errsvc"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterErrors(t *testing.T) {
	mgr := errsvc.NewManager()

	t.Run("register map should pass", func(t *testing.T) {
		mgr.MustRegisterMap(map[int32]string{503999: "test1", 403999: "none"})

		e := mgr.NewError(503999, "test2")
		assert.Equal(t, "test1", e.Message)
	})

	t.Run("register map again should panic", func(t *testing.T) {
		defer func() {
			assert.NotNil(t, recover())
		}()
		mgr.MustRegisterMap(map[int32]string{503999: "test1"})
	})

	t.Run("register < 400 should panic", func(t *testing.T) {
		defer func() {
			assert.NotNil(t, recover())
		}()
		mgr.MustRegisterMap(map[int32]string{1: "test1"})
	})
}
