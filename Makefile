## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^//'

.PHONY: confirm
confirm:
	@echo 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

## install: install dependencies
.PHONY: install
install:
	@make install/go

## install/go: install go dependencies
.PHONY: install/go
install/go:
	@go mod tidy

## test: run tests
.PHONY: test
test:
	@go test -v ./...

## lint/install: install linters
.PHONY: lint/install
lint/install:
	@echo 'Installing linters...'
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

## lint: run linters
.PHONY: lint
lint:
	@golangci-lint run --path-mode=abs --config=".golangci.yml" --timeout=5m

## lint/fix
.PHONY: lint/fix
lint/fix:
	@golangci-lint run --path-mode=abs --config=".golangci.yml" --timeout=5m --fix

## fmt
.PHONY: fmt
fmt:
	@golangci-lint fmt
