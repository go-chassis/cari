package rbac

import (
	"errors"
	"github.com/karlseguin/ccache/v2"
	"time"
)

const (
	DefaultTTL = 1 * time.Hour
)

var ErrInvalidPerms = errors.New("perms is invalid")
var ErrEmptyPerms = errors.New("perms is empty")

type FindPerms func() ([]*Role, error)
type PersistPerms func(r *Role) error

var permsCache = ccache.New(ccache.Configure())

// WritePerms save cache
func WritePerms(r *Role) error {
	permsCache.Set(r.Name, r, DefaultTTL)
	return nil
}

// ReadPerms only return data in cache
func ReadPerms(roleName string) ([]*Permission, error) {
	item := permsCache.Get(roleName)
	if item == nil || item.Value() == nil {
		return nil, ErrEmptyPerms
	}
	r, ok := item.Value().(*Role)
	if !ok {
		return nil, ErrInvalidPerms
	}
	return r.Perms, nil
}

func WriteRoles(roles []*Role) {
	for _, r := range roles {
		permsCache.Set(r.Name, r, 24*time.Hour)
	}
}
