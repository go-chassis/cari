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

package rbac

import (
	"strings"
	"sync"
)

//as a user of a backend service, he only understands resource of this service,
//to decouple authorization code from business code,
//a middleware should handle all the authorization logic, and this middleware only understand rest API,
//a resource mapping helps to maintain relations between api and resource.
var resourceMap = sync.Map{}

//PartialMap saves api partial matching
var PartialMap = map[string]string{}

// GetResource try to find resource by API path, it has preheat mechanism after program start up
// an API pattern is like /resource/:id/, /resource/{id}/,
// MUST NOT pass exact resource id to this API like /resource/100, otherwise you are facing massive memory footprint
func GetResource(apiPattern string) string {
	r, ok := resourceMap.Load(apiPattern)
	if ok {
		return r.(string)
	}
	for partialAPI, resource := range PartialMap {
		if strings.Contains(apiPattern, partialAPI) {
			resourceMap.Store(apiPattern, resource)
			return resource
		}
	}
	return ""
}

// MapResource saves the mapping from api to resource, it must be exactly match
func MapResource(api, resource string) {
	resourceMap.Store(api, resource)
}

// PartialMapResource saves the mapping from api to resource, it is partial match
func PartialMapResource(api, resource string) {
	PartialMap[api] = resource
}
