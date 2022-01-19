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

package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-chassis/foundation/gopool"
	"github.com/go-chassis/openlog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	"github.com/go-chassis/cari/db"
	"github.com/go-chassis/cari/db/config"
	"github.com/go-chassis/cari/db/mongo/log"
)

const (
	MongoCheckDelay     = 2 * time.Second
	HeathChekRetryTimes = 3
	DefaultDBName       = "servicecomb"
)

var (
	client *Client

	ErrOpenDbFailed  = errors.New("open db failed")
	ErrRootCAMissing = errors.New("rootCAFile is empty in config file")
	ErrURIIsEmpty    = errors.New("uri is empty")
)

func init() {
	db.Install("mongo", NewDatasource)
}

func NewDatasource(c *config.Config) error {
	return initClient(c)
}

func initClient(c *config.Config) error {
	NewMongoClient(c)
	select {
	case err := <-GetClient().Err():
		return err
	case <-GetClient().Ready():
		return nil
	}
}

type Client struct {
	client *mongo.Client
	db     *mongo.Database
	config *config.Config

	err       chan error
	ready     chan struct{}
	goroutine *gopool.Pool
}

func NewMongoClient(config *config.Config) {
	inst := &Client{}
	if err := inst.Initialize(config); err != nil {
		log.GetLogger().Error("failed to init mongodb" + err.Error())
		inst.err <- err
	}
	client = inst
}

func (mc *Client) Err() <-chan error {
	return mc.err
}

func (mc *Client) Ready() <-chan struct{} {
	return mc.ready
}

func (mc *Client) Close() {
	if mc.client != nil {
		if err := mc.client.Disconnect(context.TODO()); err != nil {
			log.GetLogger().Error("[close mongo client] failed disconnect the mongo client" + err.Error())
		}
	}
}

func (mc *Client) Initialize(config *config.Config) error {
	if config.Logger == nil {
		config.Logger = openlog.GetLogger()
	}
	log.SetLogger(config.Logger)
	mc.err = make(chan error, 1)
	mc.ready = make(chan struct{})
	mc.goroutine = gopool.New()
	mc.config = config
	if len(config.URI) == 0 {
		return ErrURIIsEmpty
	}
	cs, err := connstring.ParseAndValidate(config.URI)
	if err != nil {
		return err
	}
	dbName := DefaultDBName
	if len(cs.Database) != 0 {
		dbName = cs.Database
	}
	err = mc.newClient(context.Background(), dbName)
	if err != nil {
		return err
	}
	mc.startHealthCheck()
	close(mc.ready)
	return nil
}

func (mc *Client) newClient(ctx context.Context, dbName string) (err error) {
	clientOptions := []*options.ClientOptions{options.Client().ApplyURI(mc.config.URI)}
	clientOptions = append(clientOptions, options.Client().SetMaxPoolSize(uint64(mc.config.PoolSize)))
	if mc.config.SSLEnabled {
		clientOptions = append(clientOptions, options.Client().SetTLSConfig(mc.config.TLSConfig))
		log.GetLogger().Info("enabled ssl communication to mongodb")
	}
	mc.client, err = mongo.Connect(ctx, clientOptions...)
	if err != nil {
		log.GetLogger().Error("failed to connect to mongo" + err.Error())
		if derr := mc.client.Disconnect(ctx); derr != nil {
			log.GetLogger().Error("[init mongo client] failed to disconnect mongo clients" + derr.Error())
		}
		return
	}
	mc.db = mc.client.Database(dbName)
	if mc.db == nil {
		return ErrOpenDbFailed
	}
	return nil
}

func (mc *Client) startHealthCheck() {
	mc.goroutine.Do(mc.HealthCheck)
}

func (mc *Client) HealthCheck(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			mc.Close()
			return
		case <-time.After(MongoCheckDelay):
			for i := 0; i < HeathChekRetryTimes; i++ {
				err := mc.client.Ping(context.Background(), nil)
				if err == nil {
					break
				}
				log.GetLogger().Error(fmt.Sprintf("retry to connect to mongodb %s after %s",
					mc.config.URI, MongoCheckDelay) + err.Error())
				select {
				case <-ctx.Done():
					mc.Close()
					return
				case <-time.After(MongoCheckDelay):
				}
			}
		}
	}
}

func GetClient() *Client {
	return client
}

// ExecTxn execute a transaction command
// want to abort transaction, return error in cmd fn impl, otherwise it will commit transaction
func (mc *Client) ExecTxn(ctx context.Context, cmd func(sessionContext mongo.SessionContext) error) error {
	session, err := mc.client.StartSession()
	if err != nil {
		return err
	}
	if err = session.StartTransaction(); err != nil {
		return err
	}
	defer session.EndSession(ctx)
	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err = cmd(sc); err != nil {
			if err = session.AbortTransaction(sc); err != nil {
				return err
			}
		} else {
			if err = session.CommitTransaction(sc); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (mc *Client) GetDB() *mongo.Database {
	return mc.db
}

func (mc *Client) CreateIndexes(ctx context.Context, Table string, indexes []mongo.IndexModel) error {
	_, err := mc.db.Collection(Table).Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return err
	}
	return nil
}

func EnsureCollection(col string, validator interface{}, indexes []mongo.IndexModel) {
	err := client.GetDB().CreateCollection(context.Background(), col, options.CreateCollection().SetValidator(validator))
	wrapCreateCollectionError(err)

	err = client.CreateIndexes(context.Background(), col, indexes)
	wrapCreateIndexesError(err)
}

func wrapCreateCollectionError(err error) {
	if err != nil {
		if IsCollectionsExist(err) {
			log.GetLogger().Warn("collection already exist")
			return
		}
		log.GetLogger().Fatal("failed to create collection with validation")
	}
}

func wrapCreateIndexesError(err error) {
	if err != nil {
		if IsDuplicateKey(err) {
			log.GetLogger().Warn("indexes already exist")
			return
		}
		log.GetLogger().Fatal("failed to create indexes")
	}
}
