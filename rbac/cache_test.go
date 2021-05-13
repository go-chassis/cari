package rbac_test

import (
	"github.com/go-chassis/cari/rbac"
	"github.com/stretchr/testify/assert"
	"testing"
)

var roles = []*rbac.Role{
	{
		Name: "tester",
		Perms: []*rbac.Permission{{
			Resources: []*rbac.Resource{{Type: "service"}},
			Verbs:     []string{"get"},
		}},
	},
	{
		Name: "admin",
		Perms: []*rbac.Permission{{
			Resources: []*rbac.Resource{{Type: "service"}, {Type: "account"}},
			Verbs:     []string{"*"},
		}},
	},
}

func TestReadPerms(t *testing.T) {
	rbac.WriteRoles(roles)
	t.Run("given tester role, should able to get service", func(t *testing.T) {
		perms, err := rbac.ReadPerms("tester")
		assert.NoError(t, err)
		assert.Equal(t, "service", perms[0].Resources[0].Type)
		assert.Equal(t, "get", perms[0].Verbs[0])
	})
	t.Run("given tester role, after write new perms should able to create service", func(t *testing.T) {
		r := &rbac.Role{Name: "tester", Perms: []*rbac.Permission{{
			Resources: []*rbac.Resource{{Type: "service"}},
			Verbs:     []string{"get", "create"},
		}}}
		rbac.WritePerms(r)
		perms, err := rbac.ReadPerms("tester")
		assert.NoError(t, err)
		assert.Equal(t, "service", perms[0].Resources[0].Type)
		assert.Equal(t, "get", perms[0].Verbs[0])
		assert.Equal(t, "create", perms[0].Verbs[1])
	})
	t.Run("given wrong role, should return err", func(t *testing.T) {
		_, err := rbac.ReadPerms("a")
		assert.Error(t, err)
	})
}
