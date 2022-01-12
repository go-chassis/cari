# CARI
the full name is cloud application runtime interface, 
it defines the interfaces of cloud services that application should follow to run in cloud


# packages
## config
defines a common struct to store kv configuration

## discovery
defines a standard service for register and discovery, service like eureka,consul and service center 

## security
defines a common interface to encrypt and decrypt info

## rbac
defines a common struct to account and role info

### db
defines a config struct to connect db, contains two connection methods of db (etcd and mongo)

**how to use it**

1. First you should import the bootstrap package
```go
import (
    _ "github.com/go-chassis/cari/db/bootstrap"
)
```
2. The second initialization db selects etcd or mongo
```go
     // etcd
    cfg:=config.Config{
        Kind:    "mongo",
        URI:     "mongodb://127.0.0.1:27017",
        Timeout: 10 * time.Second,
    }
    // or mongo
    cfg:=config.Config{
        Kind:    "mongo",
        URI:     "mongodb://127.0.0.1:27017",
        Timeout: 10 * time.Second,
    }
    err = db.Init(&cfg)
```
3. Using etcd client
```go
    etcdadpt.Put()
```

4. Using mongo client
```go
    mongo.GetClient().GetDB()
```

## sync
defines common structs for synchronization mechanism, synchronize data to different peer clusters

## dlock
dlock provide distributed lock function

**how to use it**

1. First you should initialize db

2. Second import dlock bootstrap package
```go
    import (
        _ "github.com/go-chassis/cari/dlock/bootstrap"
    )
```

3. Init dlock
```go
    dlock.Init(dlock.Options{Kind: "etcd")
```



