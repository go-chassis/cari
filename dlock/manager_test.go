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

package dlock_test

import (
	"testing"

	"github.com/go-chassis/cari/dlock"
	"github.com/stretchr/testify/assert"

	_ "github.com/go-chassis/cari/dlock/bootstrap"
)

func TestInit(t *testing.T) {
	t.Run("init dlock should pass", func(t *testing.T) {
		err := dlock.Init(dlock.Options{
			Kind: "etcd",
		})
		assert.Nil(t, err)

		err = dlock.Init(dlock.Options{
			Kind: "mongo",
		})
		assert.Nil(t, err)
	})
}
