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

package rbac

import (
	"errors"
	"github.com/go-chassis/cari/discovery"
	"github.com/go-chassis/cari/pkg/errsvc"
)

var (
	ErrInvalidHeader      = errors.New("invalid auth header")
	ErrSameAsName         = errors.New("account name and password MUST NOT be same")
	ErrSameAsReversedName = errors.New("password MUST NOT be the revered account name")
	ErrNoHeader           = errors.New("should provide Authorization header")
	ErrInvalidCtx         = errors.New("invalid context")
	ErrConvert            = errors.New("type convert error")
	MsgConvertErr         = "type convert error"
	ErrConvertErr         = errors.New(MsgConvertErr)
)

// error code range: ***200 - ***249
const (
	ErrAccountNotExist       int32 = 400200
	ErrRoleNotExist          int32 = 400201
	ErrAccountHasInvalidRole int32 = 400202
	ErrAccountNoQuota        int32 = 400203
	ErrRoleNoQuota           int32 = 400204
	ErrRoleIsBound           int32 = 400205

	ErrUnauthorized             int32 = 401201
	ErrUserOrPwdWrong           int32 = 401202
	ErrNoPermission             int32 = 401203
	ErrNoAuthHeader             int32 = 401204
	ErrTokenExpired             int32 = 401205
	ErrTokenOwnedAccountDeleted int32 = 401206

	ErrAccountBlocked int32 = 403201

	ErrAccountConflict int32 = 409200
	ErrRoleConflict    int32 = 409201
)

var errorsMap = map[int32]string{
	ErrAccountNotExist:       "Account not exists",
	ErrRoleNotExist:          "Role not exists",
	ErrAccountHasInvalidRole: "Account has invalid role(s)",
	ErrAccountNoQuota:        "No quota to create account",
	ErrRoleNoQuota:           "No quota to create role",
	ErrRoleIsBound:           "Role is bound to some user(s)",

	ErrAccountBlocked: "Account blocked",

	ErrUnauthorized:             "Request unauthorized",
	ErrUserOrPwdWrong:           "User name or password is wrong",
	ErrNoPermission:             "No permission(s)",
	ErrNoAuthHeader:             "No authorization header",
	ErrTokenExpired:             "Token is expired",
	ErrTokenOwnedAccountDeleted: "The account that owns the token is deleted",

	ErrAccountConflict: "account name is duplicated",
	ErrRoleConflict:    "role name is duplicated",
}

func init() {
	discovery.MustRegisterErrs(errorsMap)
}

func MustRegisterErrs(errs map[int32]string) {
	discovery.MustRegisterErrs(errs)
}

func MustRegisterErr(code int32, message string) {
	discovery.MustRegisterErr(code, message)
}

func NewError(code int32, detail string) *errsvc.Error {
	return discovery.NewError(code, detail)
}
