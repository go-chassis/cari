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

package sync

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

const (
	PendingStatus = "pending"
	DoneStatus    = "done"
)

// NewTask return task with domain, project , action , resourceType and resource
func NewTask(domain, project, action, resourceType string, resource interface{}) (*Task, error) {
	taskId, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	resourceValue, err := json.Marshal(resource)
	if err != nil {
		return nil, err
	}
	return &Task{
		ID:           taskId.String(),
		Domain:       domain,
		Project:      project,
		ResourceType: resourceType,
		Resource:     resourceValue,
		Action:       action,
		Timestamp:    time.Now().UnixNano(),
		Status:       PendingStatus,
	}, nil
}
