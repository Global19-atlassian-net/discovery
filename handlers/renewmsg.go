package handlers

import "time"

// RenewMsg is used to refresh the current view of the cluster in discovery service
type RenewMsg struct {
	// ClusterID is the cluster's unique identifier
	ClusterID string `json:"clusterId"`
	// Members is a list of current cluster members. Format of each member: "name=peer-endpoint" e.g.: "machine2=http://127.0.0.17:2380"
	Members []Member `json:"members"`
	// TTL in seconds
	TTL uint64 `json:"ttl"`
	// RequestTime specifies the request time at the leader side. (for logging/debugging)
	RequestTime time.Time `json:"requestTime"`
	// LeaderName is the name of the current leader that sends this message (for logging/debugging)
	LeaderName string `json:"leaderName"`
}

// Member represents the required data for a single member to update
type Member struct {
	// ClusterID is the cluster's unique identifier
	MemberID string `json:"memberId"`
	Name     string `json:"name"`
	PeerAddr string `json:"peerAddr"`
}
