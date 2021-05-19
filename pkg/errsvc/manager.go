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

package errsvc

import (
	"fmt"
)

const initialSize = 50

type Manager struct {
	errorsMap map[int32]string
}

func (m *Manager) MustRegisterMap(errs map[int32]string) {
	for code, msg := range errs {
		m.MustRegister(code, msg)
	}
}

func (m *Manager) MustRegister(code int32, message string) {
	if code < 400000 || code >= 600000 {
		panic(fmt.Errorf("error code[%v] should be between 4xx and 5xx", code))
	}
	if _, exist := m.errorsMap[code]; exist {
		panic(fmt.Errorf("register duplicated error[%v]", code))
	}
	m.errorsMap[code] = message
}

func (m *Manager) NewError(code int32, detail string) *Error {
	return &Error{
		Code:    code,
		Message: m.errorsMap[code],
		Detail:  detail,
	}
}

func NewManager() *Manager {
	return &Manager{
		errorsMap: make(map[int32]string, initialSize),
	}
}
