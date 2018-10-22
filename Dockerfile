FROM codesimple/elm:0.18 AS frontend
RUN mkdir -p /src
WORKDIR /src
COPY elm-package.json /src/elm-package.json
COPY elm-stuff /src/elm-stuff
COPY Main.elm /src/Main.elm
RUN elm-make Main.elm --output=main.html

FROM golang:1.11.0-alpine3.8 AS backend
RUN mkdir -p /go/src/github.com/damoon/goldpinger
WORKDIR /go/src/github.com/damoon/goldpinger
COPY vendor /go/src/github.com/damoon/goldpinger/vendor
RUN go build all
COPY pkg /go/src/github.com/damoon/goldpinger/pkg
RUN go build github.com/damoon/goldpinger/pkg
COPY cmd/goldpinger /go/src/github.com/damoon/goldpinger/cmd/goldpinger
RUN go install ./cmd/goldpinger/

FROM alpine:3.8
COPY normalize.css /static/normalize.css
COPY styles.css /static/styles.css
COPY alert.png /static/alert.png
COPY --from=frontend /src/main.html /static/index.html
COPY --from=backend /go/bin/goldpinger /goldpinger
ENTRYPOINT ["/goldpinger"]
