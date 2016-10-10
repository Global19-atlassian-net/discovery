package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/coreos/etcd/client"
	ctx "golang.org/x/net/context"
)

const (
	// RenewTimeout is the TTL for each etcd entry
	// Must be refreshed by etcd leader, otherwise the entry will be removed (for the failed/removed nodes)
	RenewTimeout = 5 * 60 * time.Second

	// MaxJSONSize limits the size of input Json (1MB)
	MaxJSONSize = 1000000
)

// RenewHandler handles all PUT and GET requests related to the formation of cluster.
// TODO: make it transactional by using etcd v3 apis
func (h *Handler) RenewHandler(w http.ResponseWriter, r *http.Request) {
	// Receive list of current cluster members
	var renewMsg RenewMsg
	err := json.NewDecoder(r.Body).Decode(&renewMsg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("message received: ", renewMsg)

	clusterKey := path.Join("_etcd", "registry", renewMsg.ClusterID)
	resp, err := h.EtcdClient.Get(ctx.Background(), clusterKey, &client.GetOptions{Recursive: true})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	log.Println("resp: ", resp)

	currentMembers := renewMsg.Members
	for _, member := range currentMembers {
		log.Println("member to be refreshed or added as a part of renew: ", member)
		key := path.Join(clusterKey, member.MemberID)
		value := member.Name + "=" + member.PeerAddr
		t := time.Duration(renewMsg.TTL) * time.Second
		log.Printf("refreshing nodes: memberID: %v | member.Name-member.PeerAddr: %v | ttl: %v", key, value, t)
		resp, err := h.EtcdClient.Set(ctx.Background(), key, value, &client.SetOptions{TTL: t})
		if err != nil {
			http.Error(w, fmt.Errorf("error in renewing nodes: %v", err).Error(), 500)
			return
		}
		log.Println("resp: ", resp)
	}

	w.Write([]byte("Cluster members refreshed.\n"))
}
