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

package etcd

import (
	"sync"

	"github.com/go-chassis/etcdadpt"

	"github.com/go-chassis/cari/dlock"
)

func init() {
	dlock.Install("etcd", NewDLock)
	dlock.Install("embeded_etcd", NewDLock)
	dlock.Install("embedded_etcd", NewDLock)
}

func NewDLock() (dlock.DLock, error) {
	return &DB{lockMap: sync.Map{}}, nil
}

type DB struct {
	lockMap sync.Map
}

func (d *DB) Lock(key string, ttl int64) error {
	lock, err := etcdadpt.Lock(key, ttl)
	if err == nil {
		d.lockMap.Store(key, lock)
	}
	return err
}

func (d *DB) TryLock(key string, ttl int64) error {
	lock, err := etcdadpt.TryLock(key, ttl)
	if err == nil {
		d.lockMap.Store(key, lock)
	}
	return err
}

func (d *DB) Renew(key string) error {
	if lock, ok := d.lockMap.Load(key); ok {
		err := lock.(*etcdadpt.DLock).Refresh()
		if err != nil {
			d.lockMap.Delete(key)
		}
		return err
	}
	return dlock.ErrDLockNotExists
}

func (d *DB) IsHoldLock(key string) bool {
	if lock, ok := d.lockMap.Load(key); ok && lock != nil {
		return true
	}
	return false
}

func (d *DB) Unlock(key string) error {
	if lock, ok := d.lockMap.Load(key); ok {
		err := lock.(*etcdadpt.DLock).Unlock()
		d.lockMap.Delete(key)
		return err
	}
	return dlock.ErrDLockNotExists
}
