package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/service/cluster"
	reg "github.com/photoprism/photoprism/internal/service/cluster/registry"
)

func TestClusterNodesRegister(t *testing.T) {
	t.Run("FeatureDisabled", func(t *testing.T) {
		app, router, conf := NewApiTest()
		conf.Options().NodeType = cluster.Instance
		ClusterNodesRegister(router)

		r := PerformRequestWithBody(app, http.MethodPost, "/api/v1/cluster/nodes/register", `{"nodeName":"pp-node-01"}`)
		assert.Equal(t, http.StatusForbidden, r.Code)
	})

	t.Run("MissingToken", func(t *testing.T) {
		app, router, conf := NewApiTest()
		conf.Options().NodeType = cluster.Portal
		ClusterNodesRegister(router)

		r := PerformRequestWithBody(app, http.MethodPost, "/api/v1/cluster/nodes/register", `{"nodeName":"pp-node-01"}`)
		assert.Equal(t, http.StatusUnauthorized, r.Code)
	})

	t.Run("DriverConflict", func(t *testing.T) {
		app, router, conf := NewApiTest()
		conf.Options().NodeType = cluster.Portal
		conf.Options().PortalToken = "t0k3n"
		ClusterNodesRegister(router)

		// With SQLite driver in tests, provisioning should fail with conflict.
		r := AuthenticatedRequestWithBody(app, http.MethodPost, "/api/v1/cluster/nodes/register", `{"nodeName":"pp-node-01"}`, "t0k3n")
		assert.Equal(t, http.StatusConflict, r.Code)
		assert.Contains(t, r.Body.String(), "portal database must be MySQL/MariaDB")
	})

	t.Run("BadName", func(t *testing.T) {
		app, router, conf := NewApiTest()
		conf.Options().NodeType = cluster.Portal
		conf.Options().PortalToken = "t0k3n"
		ClusterNodesRegister(router)

		// Empty nodeName → 400
		r := AuthenticatedRequestWithBody(app, http.MethodPost, "/api/v1/cluster/nodes/register", `{"nodeName":""}`, "t0k3n")
		assert.Equal(t, http.StatusBadRequest, r.Code)
	})

	t.Run("RotateSecretPersistsDespiteDBConflict", func(t *testing.T) {
		app, router, conf := NewApiTest()
		conf.Options().NodeType = cluster.Portal
		conf.Options().PortalToken = "t0k3n"
		ClusterNodesRegister(router)

		// Pre-create node in registry so handler goes through existing-node path
		// and rotates the secret before attempting DB ensure.
		regy, err := reg.NewFileRegistry(conf)
		assert.NoError(t, err)
		n := &reg.Node{ID: "test-id", Name: "pp-node-01", Type: "instance"}
		n.Secret = "oldsecret"
		assert.NoError(t, regy.Put(n))

		r := AuthenticatedRequestWithBody(app, http.MethodPost, "/api/v1/cluster/nodes/register", `{"nodeName":"pp-node-01","rotateSecret":true}`, "t0k3n")
		assert.Equal(t, http.StatusConflict, r.Code) // DB conflict under SQLite

		// Secret should have rotated and been persisted even though DB ensure failed.
		n2, err := regy.Get("test-id")
		assert.NoError(t, err)
		assert.NotEqual(t, "oldsecret", n2.Secret)
		assert.NotEmpty(t, n2.SecretRot)
	})
}
