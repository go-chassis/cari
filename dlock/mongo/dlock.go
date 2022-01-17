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

package mongo

import (
	"errors"
	"sync"

	"github.com/go-chassis/cari/dlock"
)

func init() {
	dlock.Install("mongo", NewDLock)
}

var ErrNotImplement = errors.New("func do not implemented")

func NewDLock() (dlock.DLock, error) {
	return &DB{lockMap: sync.Map{}}, nil
}

type DB struct {
	lockMap sync.Map
}

func (d *DB) Lock(key string, ttl int64) error {
	// TODO need to implement it
	return ErrNotImplement
}

func (d *DB) TryLock(key string, ttl int64) error {
	// TODO need to implement it
	return ErrNotImplement
}

func (d *DB) Renew(key string) error {
	// TODO need to implement it
	return ErrNotImplement
}

func (d *DB) IsHoldLock(key string) bool {
	if lock, ok := d.lockMap.Load(key); ok && lock != nil {
		return true
	}
	return false
}

func (d *DB) Unlock(key string) error {
	// TODO need to implement it
	return ErrNotImplement
}
