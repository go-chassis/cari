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
	ErrHealthCheck     int32 = 400010
	ErrListKVs         int32 = 400011
	ErrRecordNotExists int32 = 400012
	ErrGetPollingData  int32 = 400013
)

var errorsMap = map[int32]string{
	ErrHealthCheck:     "Failed to check kie healthy",
	ErrListKVs:         "Failed to list Key/value",
	ErrRecordNotExists: "Micro-service version does not exist",
	ErrGetPollingData:  "Failed to get polling data",
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
