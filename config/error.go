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

package config

import (
	"github.com/go-chassis/cari/pkg/errsvc"
)

const (
	ErrHasModified         int32 = 304001
	ErrInvalidParams       int32 = 400001
	ErrHealthCheck         int32 = 400002
	ErrObserveEvent        int32 = 400003
	ErrSkipDuplicateKV     int32 = 400004
	ErrStopUpload          int32 = 400005
	ErrRequiredRecordId    int32 = 403001
	ErrRecordNotExists     int32 = 404001
	ErrRecordAlreadyExists int32 = 409001
	ErrNotEnoughQuota      int32 = 422001
	ErrInternal            int32 = 500001
)

var errorsMap = map[int32]string{
	ErrHasModified:         "kvs has modified, need to try again",
	ErrInvalidParams:       "invalid parameter(s)",
	ErrHealthCheck:         "failed to check kie healthy",
	ErrRequiredRecordId:    "required record id",
	ErrSkipDuplicateKV:     "skip overriding duplicate kvs",
	ErrStopUpload:          "stop overriding kvs after reaching the duplicate kv",
	ErrObserveEvent:        "failed to observe event",
	ErrRecordNotExists:     "record does not exist",
	ErrRecordAlreadyExists: "record already exist",
	ErrNotEnoughQuota:      "quota is not enough",
	ErrInternal:            "internal server error",
}

var errManager = errsvc.NewManager()

func init() {
	MustRegisterErrs(errorsMap)
}

func MustRegisterErrs(errs map[int32]string) {
	errManager.MustRegisterMap(errs)
}

func MustRegisterErr(code int32, message string) {
	errManager.MustRegister(code, message)
}

func NewError(code int32, detail string) *errsvc.Error {
	return errManager.NewError(code, detail)
}
