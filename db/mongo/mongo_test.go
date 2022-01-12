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

package mongo_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-chassis/cari/db/config"
	"github.com/go-chassis/cari/db/mongo"
)

func TestNewDatasource(t *testing.T) {
	t.Run("create mongo datasource should pass with no error", func(t *testing.T) {
		cfg := &config.Config{
			Kind:    "mongo",
			URI:     "mongodb://127.0.0.1:27017",
			Timeout: 10 * time.Second,
		}
		err := mongo.NewDatasource(cfg)
		assert.NoError(t, err)
	})
}

func TestCreateCollection(t *testing.T) {
	t.Run("create a abc collection should pass", func(t *testing.T) {
		cfg := &config.Config{
			Kind:    "mongo",
			URI:     "mongodb://127.0.0.1:27017",
			Timeout: 10 * time.Second,
		}
		err := mongo.NewDatasource(cfg)
		assert.NoError(t, err)
		err = mongo.GetClient().GetDB().CreateCollection(context.Background(), "abc")
		assert.NoError(t, err)
	})
}
