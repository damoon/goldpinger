
# standard lib
FROM golang:1.11.0-alpine3.8 AS goldpinger
RUN go build all
# vendor
WORKDIR /go/src/github.com/damoon/goldpinger
COPY vendor /go/src/github.com/damoon/goldpinger/vendor
RUN go build all
# package
COPY pkg /go/src/github.com/damoon/goldpinger/pkg
RUN go build github.com/damoon/goldpinger/pkg
# binary
COPY cmd/goldpinger /go/src/github.com/damoon/goldpinger/cmd/goldpinger
RUN go build -o /goldpinger ./cmd/goldpinger

# standard lib
FROM golang:1.11.0-stretch as wasm
ENV GOARCH wasm
ENV GOOS js
# https://github.com/golang/go/wiki/WebAssembly#executing-webassembly-with-nodejs
ENV PATH="$PATH:/usr/local/go/misc/wasm"
RUN curl -sL https://deb.nodesource.com/setup_9.x | bash -
RUN apt-get update && apt-get install nodejs -y
RUN go build all
# vendor
WORKDIR /go/src/github.com/damoon/goldpinger
COPY vendor /go/src/github.com/damoon/goldpinger/vendor
RUN go build ./vendor/github.com/mohae/deepcopy
# package
COPY pkg /go/src/github.com/damoon/goldpinger/pkg
RUN go test ./pkg
RUN go build github.com/damoon/goldpinger/pkg
# wasm
COPY cmd/wasm /go/src/github.com/damoon/goldpinger/cmd/wasm
RUN go test ./cmd/wasm
RUN go build -o cmd/wasm/goldpinger.wasm cmd/wasm/*.go

# development image
FROM alpine:3.8 AS dev
COPY public /public
COPY --from=wasm /go/src/github.com/damoon/goldpinger/cmd/wasm/goldpinger.wasm /public/goldpinger.wasm
COPY --from=goldpinger /goldpinger /goldpinger
ENTRYPOINT ["/goldpinger"]

# compression tools
FROM ubuntu:18.04 AS compressor
RUN apt-get update && \
    apt-get install --no-install-recommends -y zopfli brotli upx-ucl && \
    rm -rf /var/lib/apt/lists/*

# precompress static public http files
FROM compressor AS wasm-compressed
COPY --from=wasm /go/src/github.com/damoon/goldpinger/cmd/wasm/goldpinger.wasm /goldpinger.wasm
RUN brotli --best /goldpinger.wasm
RUN zopfli --i50 /goldpinger.wasm
COPY public /public
RUN cd /public && brotli --best *.css *.js *.html
RUN cd /public && zopfli --i50 *.css *.js *.html

# golang binary without debuging
FROM goldpinger AS goldpinger-prod
RUN go build -ldflags="-s -w" -o /goldpinger ./cmd/goldpinger

# compress goldpinger binary
FROM compressor AS goldpinger-compressed
COPY --from=goldpinger-prod /goldpinger /goldpinger
RUN upx --best /goldpinger

# compressed image 
FROM alpine:3.8 AS prod
COPY --from=wasm-compressed public /public
COPY --from=wasm-compressed /goldpinger.wasm /goldpinger.wasm.br /goldpinger.wasm.gz /public/
COPY --from=goldpinger-compressed /goldpinger /goldpinger
ENTRYPOINT ["/goldpinger"]