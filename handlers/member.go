package handlers

import (
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/gorilla/mux"
	ctx "golang.org/x/net/context"
)

const (
	// DefaultClientTimeout is the default timeout for etcd client
	DefaultClientTimeout = 5 * time.Second
)

// MemberHandle generates a comma-separated list of current cluster members' client endpoints.
// This list could be directly used to create new etcd client objects.
func (h *Handler) MemberHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	log.Println("token: ", token)

	key := path.Join("_etcd", "registry", token)

	resp, err := h.EtcdClient.Get(ctx.Background(), key, &client.GetOptions{Recursive: true, Quorum: true})
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	members, err := getClientEndpointsFromResp(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	str := strings.Join(members, ",")
	w.Write([]byte(str))
}

func getClientEndpointsFromResp(resp *client.Response) ([]string, error) {
	var members []string
	for _, node := range resp.Node.Nodes {
		log.Println("node: ", node)
		// We want to return a list of client endpoints that can be directly used to create etcd clients
		// However, etcd stores them as peer endpoints (which are used for communication between etcd peers).
		// The only difference is in the port number. We don't rely on predefined ports, but only assume
		// peer port number is always client-port+1
		peerEndpoint := strings.Split(node.Value, "=")[1] //format: "member-name=scheme://peer-ip:peer-port"
		u, err := url.Parse(peerEndpoint)
		if err != nil {
			return nil, err
		}
		peerSplits := strings.Split(u.Host, ":")
		peerPort, err := strconv.Atoi(peerSplits[1])
		if err != nil {
			return nil, err
		}
		clientEndpoint := url.URL{
			Scheme: u.Scheme,
			Host:   peerSplits[0] + ":" + strconv.Itoa(peerPort-1),
		}

		members = append(members, clientEndpoint.String())
	}

	return members, nil
}
