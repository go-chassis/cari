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
	"github.com/little-cui/etcdadpt"
	// support embedded etcd
	_ "github.com/little-cui/etcdadpt/embedded"
	_ "github.com/little-cui/etcdadpt/remote"

	"github.com/go-chassis/cari/db"
	"github.com/go-chassis/cari/db/config"
)

func init() {
	db.Install("etcd", NewDatasource)
	db.Install("embeded_etcd", NewDatasource)
	db.Install("embedded_etcd", NewDatasource)
}

func NewDatasource(c *config.Config) error {
	return etcdadpt.Init(etcdadpt.Config{
		Kind:             c.Kind,
		ClusterAddresses: c.URI,
		SslEnabled:       c.SSLEnabled,
		TLSConfig:        c.TLSConfig,
		Logger:           c.Logger,
	})
}
