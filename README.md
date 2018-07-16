# Goldpinger

Goldpinger checks the connection between pods in kubernetes.

## Usage

run `make` for usage

## Access

run `make proxy` to forward apiserver to port `http://127.0.0.1:8001`

user interface: http://localhost:8001/api/v1/namespaces/goldpinger/services/goldpinger/proxy/

## Update dependencies

Follow https://github.com/kubernetes/client-go/blob/master/INSTALL.md#godep to update the kubernetes client