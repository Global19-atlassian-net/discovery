package handlers

import "github.com/coreos/etcd/client"

type Handler struct {
	EtcdClient client.KeysAPI
}

func NewHandler(etcdClient client.KeysAPI) *Handler {
	return &Handler{
		EtcdClient: etcdClient,
	}
}
