# Define Go command, which can be overridden
GO ?= go

include Makefile.local

# Default target: formats code, runs the linter, and builds both agent and server binaries
.PHONY: all
all: fmt lint build-client build-server

# ----------- build Commands -----------
# build the server binary from Go source files in cmd/server directory
.PHONY: build-server

BUILDINFO_PKG_SERVER=github.com/npavlov/go-password-manager/internal/server/buildinfo
BUILD_FLAGS_SERVER=-X '$(BUILDINFO_PKG_SERVER).Version=1.0.0' \
            -X '$(BUILDINFO_PKG_SERVER).Date=$(shell date -u +%Y-%m-%d)' \
            -X '$(BUILDINFO_PKG_SERVER).Commit=$(shell git rev-parse HEAD)'

build-server:
	$(GO) build -gcflags="all=-N -l" -ldflags="${BUILD_FLAGS_SERVER}" -o bin/server ${CURDIR}/cmd/server/main.go

.PHONY: build-server-docker
build-server-docker:
	docker build -f "./deployment/dockerfile" .

# build the agent binary from Go source files in cmd/agent directory
.PHONY: build-client
BUILDINFO_PKG_CLIENT=github.com/npavlov/go-password-manager/internal/client/buildinfo
BUILD_FLAGS_CLIENT=-X '$(BUILDINFO_PKG_CLIENT).Version=1.0.0' \
            -X '$(BUILDINFO_PKG_CLIENT).Date=$(shell date -u +%Y-%m-%d)' \
            -X '$(BUILDINFO_PKG_CLIENT).Commit=$(shell git rev-parse HEAD)'

build-client:
	$(GO) build -gcflags="all=-N -l" -ldflags="${BUILD_FLAGS_CLIENT}" -o bin/client ${CURDIR}/cmd/client/main.go

# ----------- Test Commands -----------
# Run all tests and generate a coverage profile (coverage.out)
.PHONY: test
test:
	$(GO) test ./... -race -coverprofile=coverage.out -covermode=atomic

# View the test coverage report in HTML format
.PHONY: check-coverage
check-coverage:
	$(GO) tool cover -html=coverage.out

# ----------- Clean Command -----------
# Clean the bin directory by removing all generated binaries
.PHONY: clean
clean:
	rm -rf bin/

# ----------- Run Commands -----------
# Run the server directly from Go source files in cmd/server directory
.PHONY: run-server
run-server:
	$(GO) run -ldflags="${BUILD_FLAGS_SERVER}" ${CURDIR}/cmd/server/main.go

# ----------- Run Commands -----------
# Run the server in a docker image
.PHONY: run-docker
run-docker:
	docker compose -f ./deployment/docker-compose.yml up --build

# Run the agent directly from Go source files in cmd/agent directory
.PHONY: run-client
run-client:
	$(GO) run -ldflags="${BUILD_FLAGS_AGENT}" ${CURDIR}/cmd/client/main.go
.PHONY: debug-client
debug-client:
	dlv debug ${CURDIR}/cmd/client --headless --listen=:40000 --api-version=2 --accept-multiclient
# Run docker container with delve support
.PHONY: run-docker-debug
run-docker-debug:
	docker compose -f ./deployment/docker-compose-debug.yml up --build
# ----------- Lint and Format Commands -----------
# Run the linter (golangci-lint) on all Go files in the project
.PHONY: lint
lint:
	golangci-lint run ./...

# Run the linter and automatically fix issues
.PHONY: lint-fix
lint-fix:
	golangci-lint run ./... --fix

# Format all Go files in the project using the built-in Go formatting tool
.PHONY: fmt
fmt:
	$(GO) fmt ./...

# Format all Go files in the project using gofumpt for strict formatting rules
.PHONY: gofumpt
gofumpt:
	gofumpt -l -w .

# ----------- Dependency Management -----------
# Update all Go module dependencies
.PHONY: deps
deps:
	$(GO) get -u ./...

# ----------- Database Migration Commands -----------
# Create a new migration using Atlas
.PHONY: atlas-migration
atlas-migration:
	atlas migrate diff $(MIGRATION_NAME) --env dev

# Apply migrations using Goose
.PHONY: goose-up
goose-up:
	goose -dir migrations postgres "$(DATABASE_DSN)" up

# Rollback migrations using Goose
.PHONY: goose-down
goose-down:
	goose -dir migrations postgres "$(DATABASE_DSN)" down

# Generate proto contracts

.PHONY: buf-generate

buf-generate:
	buf generate

# Lint proto contracts

.PHONY: buf-lint

buf-lint:
	buf lint

# Cross-platform builds for client
.PHONY: build-client-all
build-client-all: build-client-linux build-client-mac build-client-windows

# Build client for Linux
.PHONY: build-client-linux
build-client-linux:
	GOOS=linux GOARCH=amd64 $(GO) build -gcflags="all=-N -l" -ldflags="${BUILD_FLAGS_CLIENT}" -o bin/client-linux ${CURDIR}/cmd/client/main.go

# Build client for Mac (darwin)
.PHONY: build-client-mac
build-client-mac:
	GOOS=darwin GOARCH=amd64 $(GO) build -gcflags="all=-N -l" -ldflags="${BUILD_FLAGS_CLIENT}" -o bin/client-mac ${CURDIR}/cmd/client/main.go

# Build client for Windows (with .exe extension)
.PHONY: build-client-windows
build-client-windows:
	GOOS=windows GOARCH=amd64 $(GO) build -gcflags="all=-N -l" -ldflags="${BUILD_FLAGS_CLIENT}" -o bin/client-windows.exe ${CURDIR}/cmd/client/main.go

.PHONY: generate-cert
generate-cert:
	@mkdir -p certs
	openssl req -x509 -newkey rsa:4096 -sha256 -days 365 -nodes \
		-keyout certs/key.pem -out certs/cert.pem \
		-subj "/CN=localhost" \
		-addext "subjectAltName=DNS:localhost,IP:127.0.0.1"