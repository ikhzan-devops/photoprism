package registry

import (
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/service/cluster"
)

// NodeOpts controls which optional fields get included in responses.
type NodeOpts struct {
	IncludeInternalURL bool
	IncludeDBMeta      bool
}

// NodeOptsForSession returns the default exposure policy for a session.
// Admin users see internalUrl and DB metadata; others get a redacted view.
func NodeOptsForSession(s *entity.Session) NodeOpts {
	if s != nil && s.GetUser() != nil && s.GetUser().IsAdmin() {
		return NodeOpts{IncludeInternalURL: true, IncludeDBMeta: true}
	}

	return NodeOpts{}
}

// BuildClusterNode builds a cluster.Node DTO from a registry.Node with redaction according to opts.
func BuildClusterNode(n Node, opts NodeOpts) cluster.Node {
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

// BuildClusterNodes creates a cluster node slice from the given registry node slice.
func BuildClusterNodes(list []Node, opts NodeOpts) []cluster.Node {
	if len(list) == 0 {
		return []cluster.Node{}
	}

	out := make([]cluster.Node, 0, len(list))

	for _, n := range list {
		out = append(out, BuildClusterNode(n, opts))
	}

	return out
}
