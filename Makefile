
include ./hack/help.mk
include ./hack/lint.mk

proxy: ##@development proxy kubernetes apiserver
	kubectl proxy

proxy-registry: ##@development proxy container image registry
	kubectl -n registry port-forward service/registry 5000

logs: ##@debug show and follow logs
	ktail -n goldpinger-development

top: ##@debug list containers and resource usage
	watch "kubectl -n goldpinger-development get po -o wide && kubectl top po -n goldpinger-development"

deploy-loop: ##@deploy deploy every time a file changes
	./hack/deploy-loop.sh

deploy: ##@deploy deploy once
	./hack/deploy.sh

undeploy: ##@deploy undeploy
	cat kubernetes.yaml | IMAGE="none" DOLLAR="$$" envsubst | kubectl delete -f -

dev-loop: ##@develop test and lint every time a file changes
	./hack/dev-loop.sh

test: ##develop run tests
	go test -race -cover ./cmd/goldpinger ./pkg/ ./pkg/k8s
	GOOS=js GOARCH=wasm go test -cover -exec="$(shell go env GOROOT)/misc/wasm/go_js_wasm_exec" ./cmd/wasm/

public/normalize.css:
	curl -o public/normalize.css https://necolas.github.io/normalize.css/8.0.0/normalize.css

public/wasm_exec.js:
	curl -o public/wasm_exec.js https://raw.githubusercontent.com/golang/go/release-branch.go1.11/misc/wasm/wasm_exec.js
