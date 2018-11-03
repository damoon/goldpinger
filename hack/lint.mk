
GOLINT=golangci-lint run --enable-all --fast

.PHONY: lint
lint: ##@linting Runs all linters.
	@$(MAKE) lint-bash
	@$(MAKE) lint-go

.PHONY: lint-bash
lint-bash: ##@linting Lint Bash scripts.
	shellcheck hack/*.sh

.PHONY: lint-go
lint-go: ##@linting Lint Go code.
	$(GOLINT) ./cmd/goldpinger ./pkg/ ./pkg/k8s
	GOOS=js GOARCH=wasm $(GOLINT) ./cmd/wasm
