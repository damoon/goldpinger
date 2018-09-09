
include ./hack/help.mk
include ./hack/lint.mk

proxy-api:
	kubectl proxy

proxy-registry:
	kubectl -n registry port-forward service/registry 5000

proxy:
	kubectl -n goldpinger port-forward service/goldpinger 8080:80

logs:
	ktail -n goldpinger

top:
	watch kubectl top po -n goldpinger

deploy-loop:
	CompileDaemon -pattern "(.+\\.go|.+\\.elm|.+\\.css|.+\\.yaml|.+\\.yml)$\" -build="make deploy"

deploy:
	./hack/deploy.sh

normalize.css:
	curl -o normalize.css https://necolas.github.io/normalize.css/8.0.0/normalize.css
