package api

import (
	entitypkg "github.com/photoprism/photoprism/internal/entity"
	reg "github.com/photoprism/photoprism/internal/service/cluster/registry"
	"github.com/photoprism/photoprism/pkg/service/cluster"
)

// ClusterNodeOptions controls which optional fields get included in responses.
type ClusterNodeOptions struct {
	IncludeInternalURL bool
	IncludeDBMeta      bool
}

// ClusterNodeOptionsForSession returns the default exposure policy for a session.
// Admin users see internalUrl and DB metadata; others get a redacted view.
func ClusterNodeOptionsForSession(s *entitypkg.Session) ClusterNodeOptions {
	if s != nil && s.User() != nil && s.User().IsAdmin() {
		return ClusterNodeOptions{IncludeInternalURL: true, IncludeDBMeta: true}
	}

	return ClusterNodeOptions{}
}

// BuildClusterNode builds a ClusterNode DTO from a registry.Node with redaction according to opts.
func BuildClusterNode(n reg.Node, opts ClusterNodeOptions) cluster.Node {
	out := cluster.Node{
		ID:        n.ID,
		Name:      n.Name,
		Type:      n.Type,
		Labels:    n.Labels,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}

	if opts.IncludeInternalURL && n.Internal != "" {
		out.InternalURL = n.Internal
	}

	if opts.IncludeDBMeta {
		out.DB = &cluster.NodeDB{
			Name:            n.DB.Name,
			User:            n.DB.User,
			DBLastRotatedAt: n.DB.RotAt,
		}
	}

	return out
}

// BuildClusterNodes builds DTOs for a slice of registry.Node.
func BuildClusterNodes(list []reg.Node, opts ClusterNodeOptions) []cluster.Node {
	if len(list) == 0 {
		return []cluster.Node{}
	}

	out := make([]cluster.Node, 0, len(list))

	for _, n := range list {
		out = append(out, BuildClusterNode(n, opts))
	}

	return out
}
