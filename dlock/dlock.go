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

// Package dlock provide distributed lock function
package dlock

import (
	"errors"
)

var ErrDLockNotExists = errors.New("DLock do not exist")

type DLock interface {
	Lock(key string, ttl int64) error
	TryLock(key string, ttl int64) error
	Renew(key string) error
	IsHoldLock(key string) bool
	Unlock(key string) error
}

func Lock(key string, ttl int64) error {
	return Instance().Lock(key, ttl)
}

func TryLock(key string, ttl int64) error {
	return Instance().TryLock(key, ttl)
}

func Renew(key string) error {
	return Instance().Renew(key)
}

func IsHoldLock(key string) bool {
	return Instance().IsHoldLock(key)
}

func Unlock(key string) error {
	return Instance().Unlock(key)
}
