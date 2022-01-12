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

package db

import (
	"errors"
	"fmt"

	"github.com/go-chassis/cari/db/config"
)

var (
	isInitialized bool

	plugins = make(map[string]dbInitFunc)

	ErrIsInitialized = errors.New("instance is initialized")
)

type dbInitFunc func(c *config.Config) error

func Install(pluginImplName string, f dbInitFunc) {
	plugins[pluginImplName] = f
}

func Init(c *config.Config) (err error) {
	if isInitialized {
		return nil
	}
	if f, ok := plugins[c.Kind]; ok {
		err = f(c)
	} else {
		return fmt.Errorf("this %s db type is not supported", c.Kind)
	}
	if err == nil {
		isInitialized = true
	}
	return err
}
