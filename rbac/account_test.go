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

import "testing"

func TestAccount_Check(t *testing.T) {
	type fields struct {
		ID                  string
		Name                string
		Password            string
		Role                string
		Roles               []string
		TokenExpirationTime string
		CurrentPassword     string
		Status              string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "given same name and pwd",
			fields: fields{
				Name:     "test",
				Password: "test",
			},
			wantErr: true,
		},
		{name: "given diff name and pwd",
			fields: fields{
				Name:     "test",
				Password: "Test-a1",
			},
			wantErr: false,
		},
		{name: "given reversed name as pwd",
			fields: fields{
				Name:     "test-a",
				Password: "a-tset",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		a := &Account{
			ID:                  tt.fields.ID,
			Name:                tt.fields.Name,
			Password:            tt.fields.Password,
			Role:                tt.fields.Role,
			Roles:               tt.fields.Roles,
			TokenExpirationTime: tt.fields.TokenExpirationTime,
			CurrentPassword:     tt.fields.CurrentPassword,
			Status:              tt.fields.Status,
		}
		if err := a.Check(); (err != nil) != tt.wantErr {
			t.Errorf("%q. Account.Check() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
