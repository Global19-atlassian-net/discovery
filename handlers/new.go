package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/coreos/etcd/client"
	"github.com/quantum/discovery/pkg/lockstring"
	ctx "golang.org/x/net/context"
)

const DefaultLeaderHost = "127.0.0.1:2379"

var (
	cfg           *client.Config
	discHost      string
	currentLeader lockstring.LockString
	etcdMachines  []string
)

func init() {
	currentLeader.Set(DefaultLeaderHost)
}

func generateCluster() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}

func Setup(host string) error {
	u, err := url.Parse(host)
	if err != nil {
		return err
	}
	currentLeader.Set(u.Host)
	discHost = host
	return nil
}

func (h *Handler) setupToken(size int) (string, error) {
	token := generateCluster()
	if token == "" {
		return "", errors.New("Couldn't generate a token")
	}

	key := path.Join("_etcd", "registry", token)

	resp, err := h.EtcdClient.Create(ctx.Background(), path.Join(key, "_config", "size"), strconv.Itoa(size))
	if err != nil {
		return "", fmt.Errorf("Couldn't setup state. resp: %+v | err: %+v", resp, err)
	}

	return token, nil
}

func (h *Handler) deleteToken(token string) error {
	if token == "" {
		return errors.New("No token given")
	}

	_, err := h.EtcdClient.Delete(
		ctx.Background(),
		path.Join("_etcd", "registry", token),
		&client.DeleteOptions{Recursive: true},
	)

	return err
}

func (h *Handler) NewTokenHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	size := 3
	s := r.FormValue("size")
	if s != "" {
		size, err = strconv.Atoi(s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	token, err := h.setupToken(size)

	if err != nil {
		log.Printf("setupToken returned: %v", err)
		http.Error(w, "Unable to generate token", 400)
		return
	}

	log.Println("New cluster created", token)
	fmt.Fprintf(w, "%s/%s", bytes.TrimRight([]byte(discHost), "/"), token)
}
