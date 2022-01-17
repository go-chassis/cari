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

package tls

import (
	"crypto/tls"
	"errors"
	"io/ioutil"

	"github.com/go-chassis/cari/db/config"
	"github.com/go-chassis/foundation/stringutil"
	"github.com/go-chassis/foundation/tlsutil"
	"github.com/go-chassis/go-chassis/v2/security/cipher"
)

var ErrRootCAMissing = errors.New("rootCAFile is empty in config file")

func NewTLSConfig(c *config.Config) (*tls.Config, error) {
	var password string
	if c.CertPwdFile != "" {
		pwdBytes, err := ioutil.ReadFile(c.CertPwdFile)
		if err != nil {
			return nil, err
		}
		password = TryDecrypt(stringutil.Bytes2str(pwdBytes))
	}
	if c.RootCA == "" {
		return nil, ErrRootCAMissing
	}
	opts := append(tlsutil.DefaultClientTLSOptions(),
		tlsutil.WithVerifyPeer(c.VerifyPeer),
		tlsutil.WithVerifyHostName(false),
		tlsutil.WithKeyPass(password),
		tlsutil.WithCA(c.RootCA),
		tlsutil.WithCert(c.CertFile),
		tlsutil.WithKey(c.KeyFile),
	)
	return tlsutil.GetClientTLSConfig(opts...)
}

// TryDecrypt return the src when decrypt failed
func TryDecrypt(src string) string {
	res, err := cipher.Decrypt(src)
	if err != nil {
		res = src
	}
	return res
}
