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

package config

//KVDoc is database struct to store kv
type KVDoc struct {
	ID             string `json:"id,omitempty" bson:"id,omitempty" yaml:"id,omitempty" swag:"string"`
	LabelFormat    string `json:"label_format,omitempty" bson:"label_format,omitempty" yaml:"label_format,omitempty"`
	Key            string `json:"key" yaml:"key" validate:"min=1,max=128,key"`
	Value          string `json:"value" yaml:"value" validate:"max=131072,value"`                                                    //128K
	ValueType      string `json:"value_type,omitempty" bson:"value_type,omitempty" yaml:"value_type,omitempty" validate:"valueType"` //ini,json,text,yaml,properties,xml
	Checker        string `json:"check,omitempty" yaml:"check,omitempty" validate:"max=1048576,check"`                               //python script
	CreateRevision int64  `json:"create_revision,omitempty" bson:"create_revision," yaml:"create_revision,omitempty"`
	UpdateRevision int64  `json:"update_revision,omitempty" bson:"update_revision," yaml:"update_revision,omitempty"`
	Project        string `json:"project,omitempty" yaml:"project,omitempty" validate:"min=1,max=256,commonName"`
	Status         string `json:"status,omitempty" yaml:"status,omitempty" validate:"kvStatus"`
	CreateTime     int64  `json:"create_time,omitempty" bson:"create_time," yaml:"create_time,omitempty"`
	UpdateTime     int64  `json:"update_time,omitempty" bson:"update_time," yaml:"update_time,omitempty"`

	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty" validate:"max=6,dive,keys,labelK,endkeys,labelV"` //redundant
	Domain string            `json:"domain,omitempty" yaml:"domain,omitempty" validate:"min=1,max=256,commonName"`              //redundant
}
