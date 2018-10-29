
FROM golang:1.11.0-alpine3.8 AS workbench
RUN go build all
RUN GOARCH=wasm GOOS=js go build all

RUN mkdir -p /go/src/github.com/damoon/goldpinger
WORKDIR /go/src/github.com/damoon/goldpinger
COPY vendor /go/src/github.com/damoon/goldpinger/vendor
RUN go build all
RUN GOARCH=wasm GOOS=js go build ./vendor/github.com/mohae/deepcopy

COPY pkg /go/src/github.com/damoon/goldpinger/pkg
RUN go build github.com/damoon/goldpinger/pkg
RUN GOARCH=wasm GOOS=js go build github.com/damoon/goldpinger/pkg

FROM workbench AS goldpinger
COPY cmd/goldpinger /go/src/github.com/damoon/goldpinger/cmd/goldpinger
RUN go build -o /goldpinger ./cmd/goldpinger

FROM workbench as wasm
COPY cmd/wasm /go/src/github.com/damoon/goldpinger/cmd/wasm
RUN GOARCH=wasm GOOS=js go build -o cmd/wasm/goldpinger.wasm cmd/wasm/*.go

FROM alpine:3.8 AS dev
COPY public /public
COPY --from=wasm /go/src/github.com/damoon/goldpinger/cmd/wasm/goldpinger.wasm /public/goldpinger.wasm
COPY --from=goldpinger /goldpinger /goldpinger
ENTRYPOINT ["/goldpinger"]

FROM workbench AS goldpinger-prod
COPY cmd/goldpinger /go/src/github.com/damoon/goldpinger/cmd/goldpinger
RUN go build -ldflags="-s -w" -o /goldpinger ./cmd/goldpinger

FROM ubuntu:18.04 AS compressor
RUN apt-get update && \
    apt-get install --no-install-recommends -y zopfli brotli upx-ucl && \
    rm -rf /var/lib/apt/lists/*

FROM compressor AS wasm-compressed
COPY --from=wasm /go/src/github.com/damoon/goldpinger/cmd/wasm/goldpinger.wasm /goldpinger.wasm
RUN brotli --best /goldpinger.wasm
RUN zopfli --i50 /goldpinger.wasm
COPY public /public
RUN cd /public && brotli --best *.css *.js *.html
RUN cd /public && zopfli --i50 *.css *.js *.html

FROM compressor AS goldpinger-compressed
COPY --from=goldpinger-prod /goldpinger /goldpinger
RUN upx --best /goldpinger

FROM alpine:3.8 AS prod
COPY --from=wasm-compressed public /public
COPY --from=wasm-compressed /goldpinger.wasm /goldpinger.wasm.br /goldpinger.wasm.gz /public/
COPY --from=goldpinger-compressed /goldpinger /goldpinger
ENTRYPOINT ["/goldpinger"]