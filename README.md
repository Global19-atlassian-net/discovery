# Discovery Service

The `discovery.etcd.io` package is only usable during the initial bootstrap of the etcd cluster. We extended this package to support for dynamic maintenance of
cluster membership while keeping it backward-compatible with the original package.


# APIs
- `GET /new?size=N` : returns a clsuter ID for bootstrapping a cluster of size *N*
- `GET /members/CLUSTER_TOKEN` : returns a comma-separated list of *latest* members of the cluster.
- `POST /renew`: this API is called by a cluster manager periodically to update the etcd cluster's latest quorum members.
- `GET /health` : returns the health of the discovdery service. the output is *OK* or an error message.
- `GET /CLUSTER_TOKEN/{machine}` : this is used by etcd instances to form a quorum during the bootstrap process.

# Configuration

The service has three configuration options, and can be configured with runtime arguments.

* `--etcd-urls` / `DISC_ETCD_URLS`: comma separated list of backend etcd cluster members. Default value is set to `"http://127.0.0.1:2379"`
* `--addr` / `DISC_ADDR`: address:port of the discovery service.


## Docker Container

You may run the service in a docker container:

```
docker pull quantum/discovery
docker run -d -p 80:8087 -e DISC_ADDR=... -e DISC_ETCD_URLS=... quantum/discovery
```


## Development

The discovery service uses gin for easy development. It is simple to get started:

install: `go get github.com/codegangsta/gin`

```
cd cmd
gin -appPort 8087
curl --verbose -X PUT localhost:3000/new
```

# Discussion

- This implementation provides security by obscurity.

