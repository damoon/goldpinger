
FROM golang:1.11.0-alpine3.8 AS backend
RUN go build all
RUN GOARCH=wasm GOOS=js go build all

RUN mkdir -p /go/src/github.com/damoon/goldpinger
WORKDIR /go/src/github.com/damoon/goldpinger
COPY vendor /go/src/github.com/damoon/goldpinger/vendor
RUN go build all
#RUN GOARCH=wasm GOOS=js go build all

COPY pkg /go/src/github.com/damoon/goldpinger/pkg
RUN go build github.com/damoon/goldpinger/pkg
#RUN GOARCH=wasm GOOS=js go build github.com/damoon/goldpinger/pkg

COPY cmd/goldpinger /go/src/github.com/damoon/goldpinger/cmd/goldpinger
RUN go install ./cmd/goldpinger/
COPY cmd/wasm /go/src/github.com/damoon/goldpinger/cmd/wasm
RUN GOARCH=wasm GOOS=js go build -o cmd/wasm/goldpinger.wasm cmd/wasm/main.go

FROM alpine:3.8
COPY public /public
COPY --from=backend /go/src/github.com/damoon/goldpinger/cmd/wasm/goldpinger.wasm /public/goldpinger.wasm
COPY --from=backend /go/bin/goldpinger /goldpinger
ENTRYPOINT ["/goldpinger"]
