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
	"time"

	"github.com/go-chassis/go-archaius"
	"github.com/stretchr/testify/assert"

	"github.com/go-chassis/cari/db"
	_ "github.com/go-chassis/cari/db/bootstrap"
	"github.com/go-chassis/cari/db/config"
	"github.com/go-chassis/cari/dlock"
	_ "github.com/go-chassis/cari/dlock/bootstrap"
)

var (
	defaultTestDB    = "etcd"
	defaultTestDBURI = "http://127.0.0.1:2379"
)

var dbCfg = config.Config{}

func init() {
	err := archaius.Init(archaius.WithMemorySource(), archaius.WithENVSource())
	if err != nil {
		panic(err)
	}
	mode, ok := archaius.Get("TEST_DB_MODE").(string)
	if ok && mode != "" {
		defaultTestDB = mode
	}
	uri, ok := archaius.Get("TEST_DB_URI").(string)
	if ok && uri != "" {
		defaultTestDBURI = uri
	}
	dbCfg.Kind = defaultTestDB
	dbCfg.URI = defaultTestDBURI
	dbCfg.Timeout = 10 * time.Second
	err = db.Init(&dbCfg)
	if err != nil {
		panic(err)
	}
	err = dlock.Init(dlock.Options{Kind: dbCfg.Kind})
	if err != nil {
		panic(err)
	}
}

func isEtcd() bool {
	return dbCfg.Kind == "etcd"
}

func TestDLock(t *testing.T) {
	if !isEtcd() {
		return
	}
	t.Run("test lock", func(t *testing.T) {
		t.Run("lock the global key for 5s should pass", func(t *testing.T) {
			err := dlock.Lock("global", 5)
			assert.Nil(t, err)
			isHold := dlock.IsHoldLock("global")
			assert.Equal(t, true, isHold)
		})
		t.Run("two locks fight for the same lock 5s, one lock should pass, another lock should fail", func(t *testing.T) {
			err := dlock.Lock("same-lock", 5)
			assert.Nil(t, err)
			isHold := dlock.IsHoldLock("same-lock")
			assert.Equal(t, true, isHold)
			err = dlock.TryLock("same-lock", 5)
			assert.NotNil(t, err)
		})
	})
	t.Run("test try lock", func(t *testing.T) {
		t.Run("try lock the try key for 5s should pass", func(t *testing.T) {
			err := dlock.TryLock("try-lock", 5)
			assert.Nil(t, err)
			isHold := dlock.IsHoldLock("try-lock")
			assert.Equal(t, true, isHold)
			err = dlock.TryLock("try-lock", 5)
			assert.NotNil(t, err)
		})
	})
	t.Run("test renew", func(t *testing.T) {
		t.Run("renew the renew key for 5s should pass", func(t *testing.T) {
			err := dlock.Lock("renew", 5)
			assert.Nil(t, err)
			isHold := dlock.IsHoldLock("renew")
			assert.Equal(t, true, isHold)
			time.Sleep(3 * time.Second)
			err = dlock.Renew("renew")
			time.Sleep(2 * time.Second)
			err = dlock.TryLock("renew", 5)
			assert.NotNil(t, err)
		})
	})
	t.Run("test isHoldLock", func(t *testing.T) {
		t.Run("already owns the lock should pass", func(t *testing.T) {
			err := dlock.Lock("hold-lock", 5)
			assert.Nil(t, err)
			isHold := dlock.IsHoldLock("hold-lock")
			assert.Equal(t, true, isHold)
		})
		t.Run("key does not exist should fail", func(t *testing.T) {
			isHold := dlock.IsHoldLock("not-exist")
			assert.Equal(t, false, isHold)
		})
	})
	t.Run("test unlock", func(t *testing.T) {
		t.Run("unlock the unlock key should pass", func(t *testing.T) {
			err := dlock.Lock("unlock", 5)
			assert.Nil(t, err)
			err = dlock.Unlock("unlock")
			assert.Nil(t, err)
			lock := dlock.IsHoldLock("unlock")
			assert.False(t, lock)
		})
	})
}
