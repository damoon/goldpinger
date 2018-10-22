# Goldpinger

Goldpinger checks the connection between pods in kubernetes.

## Usage

run `make` for usage

## Development

Run `make proxy-registry` to connect up to your development cluster.

Run `make deploy` to deploy once. \
alertnative: \
Run `make deploy-loop` to use filewatcher for rebuilding.

Run `make top` to see eesource usage and a list of containers.

Run `make logs` to follow the logs of all containers.

## Access

Run `make proxy` to forward the kubernetes apiserver to port `http://127.0.0.1:8001`

Open http://localhost:8001/api/v1/namespaces/goldpinger-development/services/goldpinger/proxy/ to view the UI.

## Update dependencies

Follow https://github.com/kubernetes/client-go/blob/master/INSTALL.md#godep to update the kubernetes client