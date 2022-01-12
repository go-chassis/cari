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

package dlock

import (
	"errors"
	"fmt"
)

type initFunc func() (DLock, error)

var (
	plugins  = make(map[string]initFunc)
	instance DLock

	ErrKindIsEmpty   = errors.New("dlock's kind is empty")
	ErrIsInitialized = errors.New("instance has been initialized")
)

func Install(pluginImplName string, f initFunc) {
	plugins[pluginImplName] = f
}

func Init(opts Options) (err error) {
	if opts.Kind == "" {
		return ErrKindIsEmpty
	}
	if instance != nil {
		return nil
	}
	engineFunc, ok := plugins[opts.Kind]
	if !ok {
		return fmt.Errorf("plugin implement not supported [%s]", opts.Kind)
	}
	instance, err = engineFunc()
	return err
}

func Instance() DLock {
	return instance
}
