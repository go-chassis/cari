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

package discovery

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	Code    int32  `json:"errorCode,string"`
	Message string `json:"errorMessage"`
	Detail  string `json:"detail,omitempty"`
}

func (e *Error) Error() string {
	if len(e.Detail) == 0 {
		return e.Message
	}
	return e.Message + "(" + e.Detail + ")"
}

func (e *Error) Marshal() []byte {
	bs, _ := json.Marshal(e)
	return bs
}

func (e *Error) StatusCode() int {
	return int(e.Code / 1000)
}

func (e *Error) InternalError() bool {
	return e.Code >= 500000
}

func NewError(code int32, detail string) *Error {
	return &Error{
		Code:    code,
		Message: errors[code],
		Detail:  detail,
	}
}

func NewErrorf(code int32, format string, args ...interface{}) *Error {
	return NewError(code, fmt.Sprintf(format, args...))
}
