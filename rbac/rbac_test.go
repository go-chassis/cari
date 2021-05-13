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

package rbac_test

import (
	"context"
	"github.com/go-chassis/cari/rbac"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromContext(t *testing.T) {
	ctx := rbac.NewContext(context.TODO(), map[string]interface{}{
		rbac.ClaimsUser:  "root",
		rbac.ClaimsRoles: []interface{}{"admin"},
	})

	claims, _ := rbac.FromContext(ctx)
	u := claims[rbac.ClaimsUser]
	r := claims[rbac.ClaimsRoles]
	assert.Equal(t, "root", u)
	assert.Equal(t, []interface{}{"admin"}, r)

	a, err := rbac.AccountFromContext(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "root", a.Name)
	assert.Equal(t, []string{"admin"}, a.Roles)
}

func TestMustAuth(t *testing.T) {
	rbac.Add2WhiteAPIList("/test")
	rbac.Add2WhiteAPIList("/test1")
	assert.False(t, rbac.MustAuth("/test"))
	assert.False(t, rbac.MustAuth("/test1"))
	assert.True(t, rbac.MustAuth("/auth"))
	assert.True(t, rbac.MustAuth("/version"))
	assert.True(t, rbac.MustAuth("/v4/a/registry/version"))
	assert.True(t, rbac.MustAuth("/health"))
	assert.True(t, rbac.MustAuth("/v4/a/registry/health"))
}
