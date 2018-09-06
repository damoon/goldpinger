
include ./hack/help.mk
include ./hack/lint.mk

.PHONY: proxy-api
proxy-api:
	kubectl proxy

.PHONY: proxy-registry
proxy-registry:
	kubectl -n registry port-forward service/registry 5000

.PHONY: deploy-loop
deploy-loop:
	CompileDaemon -pattern "(.+\\.go|.+\\.elm|.+\\.css|.+\\.yaml|.+\\.yml)$\" -build="make deploy"

.PHONY: deploy
deploy:
	./hack/deploy.sh

normalize.css:
	curl -o normalize.css https://necolas.github.io/normalize.css/8.0.0/normalize.css
