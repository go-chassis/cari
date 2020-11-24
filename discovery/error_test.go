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
package discovery_test

import (
	"github.com/go-chassis/cari/discovery"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestError_StatusCode(t *testing.T) {
	e := discovery.Error{Code: 503999}
	if e.StatusCode() != http.StatusServiceUnavailable {
		t.Fatalf("TestError_StatusCode %v failed", e)
	}

	if !e.InternalError() {
		t.Fatalf("TestInternalError failed")
	}
}

func TestNewError(t *testing.T) {
	var err error
	err = discovery.NewError(discovery.ErrInvalidParams, "test1")
	if err == nil {
		t.Fatalf("TestNewError failed")
	}
	err = discovery.NewErrorf(discovery.ErrInvalidParams, "%s", "test2")
	if err == nil {
		t.Fatalf("TestNewErrorf failed")
	}

	if len(err.Error()) == 0 {
		t.Fatalf("TestError failed")
	}

	if len(err.(*discovery.Error).Marshal()) == 0 {
		t.Fatalf("TestMarshal failed")
	}

	if err.(*discovery.Error).StatusCode() != http.StatusBadRequest {
		t.Fatalf("TestStatusCode failed, %d", err.(*discovery.Error).StatusCode())
	}

	if err.(*discovery.Error).InternalError() {
		t.Fatalf("TestInternalError failed")
	}

	err = discovery.NewErrorf(discovery.ErrInvalidParams, "")
	if len(err.Error()) == 0 {
		t.Fatalf("TestNewErrorf with empty detial failed")
	}
}

func TestRegisterErrors(t *testing.T) {
	discovery.RegisterErrors(map[int32]string{503999: "test1", 1: "none"})

	e := discovery.NewError(503999, "test2")
	assert.Equal(t, "test1", e.Message)
}
