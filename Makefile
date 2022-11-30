export GOFLAGS=-mod=mod

GOLANG_LINTER="./.bin/golangci-lint"
GOENV:=$(shell go env GOPATH)
export PATH:=$(PATH):$(GOENV)/bin

.DEFAULT_GOAL := default
.PHONY: default
default: build

.PHONY: checkgosum
checkgosum:
	@echo "Running check go sum..."
	@test -f go.sum || { echo "go.sum missing"; exit 1; }

.PHONY: checkfmt
checkfmt:
ifneq ($(shell gofmt -l . | grep -v vendor\/),)
	@gofmt -l . | grep -v vendor\/
	@exit 1
else
	@true
endif

.PHONY: mkdirs
mkdirs:
	@mkdir ./.bin 2>/dev/null || true

.PHONY: installlint
installlint:
	@[ ! -f "./.bin/golangci-lint" ] && echo "Installing linter..." || true
	@[ ! -f "./.bin/golangci-lint" ] && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./.bin v1.50.1 >/dev/null 2>&1 || true

## lint: lints using golangci
PHONY: lint
lint:
	@ $(GOLANG_LINTER) run ./...

.PHONY: installmockgen
installmockgen:
	@which mockgen >/dev/null 2>&1 || go get github.com/golang/mock/mockgen >/dev/null 2>&1

generate: installmockgen
	@echo "Generating files..."
	@mkdir mocks 2>/dev/null || true
	@go generate ./... > /dev/null 2>&1

.PHONY: test
test:
	@echo "Running tests..."
	@CGO_ENABLED=1 go test -race -tags=unit -v ./...

.PHONY: tidy
## tidy: tidy dependencies
tidy: fetch_deps
	@echo "Tidying modules..."
	@go mod tidy

.PHONY: fetch_deps
fetch_deps:
	@echo "Fetching dependencies..."
	@go get ./...

.PHONY: clean-vendor
clean-vendor:
	@echo "Cleaning up vendors folder..."
	@rm -rf vendor

.PHONY: vendor
vendor:
	@echo "Vendoring..."
	@go mod vendor

.PHONY: pre-build
pre-build: tidy clean-vendor vendor

.PHONY: build-go-binary
build-go-binary:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -trimpath -o bin/app cmd/app/main.go

.PHONY: docker-up
docker-up:
	docker-compose -f docker-compose.yaml up --build

.PHONY: docker-down
docker-down: ## Stop docker containers and clear artefacts.
	docker-compose -f docker-compose.yaml down
#	docker system prune 

.PHONY: build
build: checkgosum  checkfmt  mkdirs installlint lint generate test pre-build build-go-binary
# docker-down docker-up