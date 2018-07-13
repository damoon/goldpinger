FROM codesimple/elm:0.18 AS frontend
RUN mkdir -p /src
WORKDIR /src
COPY elm-package.json /src/elm-package.json
COPY elm-stuff /src/elm-stuff
COPY Main.elm /src/Main.elm
RUN elm-make Main.elm --output=main.html

FROM golang:1.10.3-alpine3.8 AS backend
RUN mkdir -p /go/src/github.com/damoon/goldpinger
WORKDIR /go/src/github.com/damoon/goldpinger
COPY vendor /go/src/github.com/damoon/goldpinger/vendor
COPY pkg /go/src/github.com/damoon/goldpinger/pkg
COPY main.go /go/src/github.com/damoon/goldpinger/main.go
RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo .

FROM scratch
COPY normalize.css normalize.css
COPY styles.css /styles.css
COPY --from=frontend /src/main.html /main.html
COPY --from=backend /go/bin/goldpinger /goldpinger
ENTRYPOINT ["/goldpinger"]
