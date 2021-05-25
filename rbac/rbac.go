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

// Package rbac extract common functions to help component to implement a auth system
package rbac

import (
	"context"
	"errors"
	mapset "github.com/deckarep/golang-set"
)

const (
	ClaimsUser = "account"
	//Deprecated
	ClaimsRole  = "role"
	ClaimsRoles = "roles"

	RoleAdmin     = "admin"
	RoleDeveloper = "developer"
)

var whiteAPIList = mapset.NewSet()

func AccountFromContext(ctx context.Context) (*Account, error) {
	m, err := FromContext(ctx)
	if err != nil {
		return nil, err
	}
	accountNameI := m[ClaimsUser]
	a, ok := accountNameI.(string)
	if !ok {
		return nil, ErrConvert
	}
	roles := m[ClaimsRoles]
	roleList, err := getRolesList(roles)
	if err != nil {
		return nil, ErrConvert
	}
	account := &Account{Name: a, Roles: roleList}
	return account, nil
}

// RoleFromContext only return role name
func RolesFromContext(ctx context.Context) ([]string, error) {
	m, err := FromContext(ctx)
	if err != nil {
		return nil, err
	}
	roles := m[ClaimsRoles]
	roleList, err := getRolesList(roles)
	if err != nil {
		return nil, ErrConvert
	}
	return roleList, nil
}

func getRolesList(v interface{}) ([]string, error) {
	s, ok := v.([]interface{})
	if !ok {
		return nil, ErrConvert
	}
	rolesList := make([]string, 0)
	for _, v := range s {
		role, ok := v.(string)
		if !ok {
			return nil, ErrConvert
		}
		rolesList = append(rolesList, role)
	}
	return rolesList, nil
}

func Add2WhiteAPIList(path ...string) {
	for _, p := range path {
		whiteAPIList.Add(p)
	}
}

func MustAuth(pattern string) bool {
	return !whiteAPIList.Contains(pattern)
}

func GetRolesList(m map[string]interface{}) ([]string, error) {
	role, ok := m[ClaimsRole]
	if ok {
		r, ok := role.(string)
		if !ok {
			return nil, errors.New("convert role to string failed")
		}
		return []string{r}, nil
	}

	roles, ok := m[ClaimsRoles]
	if !ok {
		return nil, errors.New("token contains no valid roles")
	}
	roleList, err := getRolesList(roles)
	if err != nil {
		return nil, ErrConvertErr
	}
	return roleList, nil
}

//BuildResourceList join the resource to an array
func BuildResourceList(resourceType ...string) []*Resource {
	rt := make([]*Resource, len(resourceType))
	for i := 0; i < len(resourceType); i++ {
		rt[i] = &Resource{
			Type: resourceType[i],
		}
	}
	return rt
}
