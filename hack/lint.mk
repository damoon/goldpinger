
.PHONY: lint
lint: ##@linting Runs all linters.
	@$(MAKE) lint-bash
	@$(MAKE) lint-go

.PHONY: lint-bash
lint-bash: ##@linting Lint Bash scripts.
	shellcheck hack/*.sh

.PHONY: lint-go
lint-go: ##@linting Lint Go code.
	golangci-lint run
#	gometalinter --enable-all --line-length=120 -t --vendor ./...
