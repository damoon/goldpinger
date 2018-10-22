
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
	CompileDaemon -pattern "(.+\\.go|.+\\.elm|.+\\.css|.+\\.yaml|.+\\.yml)$\" -build="make deploy"
	$(MAKE) undeploy

deploy: ##@deploy deploy once
	./hack/deploy.sh

undeploy: ##@deploy undeploy
	cat kubernetes.yaml | IMAGE="none" DOLLAR="$$" envsubst | kubectl delete -f -

normalize.css:
	curl -o normalize.css https://necolas.github.io/normalize.css/8.0.0/normalize.css
