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

const (
	CreateAction = "create"
	UpdateAction = "update"
	DeleteAction = "delete"
)

// Task is db struct to store sync task
type Task struct {
	ID           string `json:"id" bson:"id"`
	Domain       string `json:"domain" bson:"domain"`
	Project      string `json:"project" bson:"project"`
	ResourceType string `json:"resource_type" bson:"resource_type"`
	Resource     []byte `json:"resource" bson:"resource"`
	Action       string `json:"action" bson:"action"`
	Timestamp    int64  `json:"timestamp" bson:"timestamp"`
	Status       string `json:"status" bson:"status"`
}

// Tombstone is db struct to store the deleted resource information
type Tombstone struct {
	ResourceID   string `json:"resource_id" bson:"resource_id"`
	ResourceType string `json:"resource_type" bson:"resource_type"`
	Domain       string `json:"domain" bson:"domain"`
	Project      string `json:"project" bson:"project"`
	Timestamp    int64  `json:"timestamp" bson:"timestamp"`
}
