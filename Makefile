
.PHONY: live-elm
live-elm:
	elm-live Main.elm

.PHONY: live-go
live-go:
	CompileDaemon -build="go build -o goldpinger main.go" -command="./goldpinger"

.PHONY: proxy
proxy:
	kubectl proxy

.PHONY: deploy-loop
deploy-loop:
	./hack/deploy-loop.sh

.PHONY: deploy
deploy:
	./hack/deploy.sh

normalize.css:
	curl -o normalize.css https://necolas.github.io/normalize.css/8.0.0/normalize.css
