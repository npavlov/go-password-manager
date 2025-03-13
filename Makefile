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

# build the agent binary from Go source files in cmd/agent directory
.PHONY: build-client
BUILDINFO_PKG_AGENT=github.com/npavlov/go-password-manager/internal/client/buildinfo
BUILD_FLAGS_AGENT=-X '$(BUILDINFO_PKG_AGENT).Version=1.0.0' \
            -X '$(BUILDINFO_PKG_AGENT).Date=$(shell date -u +%Y-%m-%d)' \
            -X '$(BUILDINFO_PKG_AGENT).Commit=$(shell git rev-parse HEAD)'

build-client:
	$(GO) build -gcflags="all=-N -l" -ldflags="${BUILD_FLAGS_AGENT}" -o bin/agent ${CURDIR}/cmd/agent/main.go

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

# Run the agent directly from Go source files in cmd/agent directory
.PHONY: run-client
run-client:
	$(GO) run -ldflags="${BUILD_FLAGS_AGENT}" ${CURDIR}/cmd/client/main.go
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