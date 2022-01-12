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

package db_test

import (
	"testing"
	"time"

	"github.com/go-chassis/go-archaius"
	"github.com/stretchr/testify/assert"

	"github.com/go-chassis/cari/db"
	_ "github.com/go-chassis/cari/db/bootstrap"
	"github.com/go-chassis/cari/db/config"
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
}

func TestInit(t *testing.T) {
	t.Run("initialize db should pass", func(t *testing.T) {
		err := db.Init(&dbCfg)
		assert.Nil(t, err)
	})
}
