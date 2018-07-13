run:
	elm-reactor --address=0.0.0.0

live-elm:
	elm-live Main.elm

live-go:
	CompileDaemon -build="go build -o goldpinger main.go" -command="./goldpinger"

.PHONY: image
image: normalize.css
	docker build .

normalize.css:
	curl -o normalize.css https://necolas.github.io/normalize.css/8.0.0/normalize.css
