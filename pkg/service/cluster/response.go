package cluster

// NodeDB represents database metadata returned for a node.
// swagger:model NodeDB
type NodeDB struct {
	Name            string `json:"name"`
	User            string `json:"user"`
	DBLastRotatedAt string `json:"dbLastRotatedAt"`
}

// Node is the API response DTO for a cluster node.
// swagger:model Node
type Node struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	InternalURL string            `json:"internalUrl,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	CreatedAt   string            `json:"createdAt"`
	UpdatedAt   string            `json:"updatedAt"`
	DB          *NodeDB           `json:"db,omitempty"`
}

// DBInfo provides basic database connection metadata for summary endpoints.
// swagger:model DBInfo
type DBInfo struct {
	Driver string `json:"driver"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
}

// SummaryResponse is the response type for GET /api/v1/cluster.
// swagger:model SummaryResponse
type SummaryResponse struct {
	PortalUUID string `json:"portalUUID"`
	Nodes      int    `json:"nodes"`
	DB         DBInfo `json:"db"`
	Time       string `json:"time"`
}

// RegisterSecrets contains newly issued or rotated node secrets.
// swagger:model RegisterSecrets
type RegisterSecrets struct {
	NodeSecret              string `json:"nodeSecret,omitempty"`
	NodeSecretLastRotatedAt string `json:"nodeSecretLastRotatedAt,omitempty"`
}

// RegisterDB describes database credentials returned during registration/rotation.
// swagger:model RegisterDB
type RegisterDB struct {
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Name            string `json:"name"`
	User            string `json:"user"`
	Password        string `json:"password,omitempty"`
	DSN             string `json:"dsn,omitempty"`
	DBLastRotatedAt string `json:"dbLastRotatedAt,omitempty"`
}

// RegisterResponse is the response body for POST /api/v1/cluster/nodes/register.
// swagger:model RegisterResponse
type RegisterResponse struct {
	Node               Node             `json:"node"`
	DB                 RegisterDB       `json:"db"`
	Secrets            *RegisterSecrets `json:"secrets,omitempty"`
	AlreadyRegistered  bool             `json:"alreadyRegistered"`
	AlreadyProvisioned bool             `json:"alreadyProvisioned"`
}

// StatusResponse is a generic status wrapper for simple ok responses.
// swagger:model StatusResponse
type StatusResponse struct {
	Status string `json:"status"`
}
