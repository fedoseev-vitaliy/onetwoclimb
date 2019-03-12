BASEPATH = $(shell pwd)

# Basic go commands
GOCMD      = go
GOBUILD    = $(GOCMD) build
GOINSTALL  = $(GOCMD) install
GORUN      = $(GOCMD) run
GOCLEAN    = $(GOCMD) clean
GOTEST     = $(GOCMD) test
GOGET      = $(GOCMD) get
GOFMT      = $(GOCMD) fmt
GOGENERATE = $(GOCMD) generate
GOTYPE     = $(GOCMD)type

# Docker
DOCKER_COMPOSE = docker-compose

# Swagger
SWAGGER = swagger

BINARY = onetwoclimb

BUILD_DIR = $(BASEPATH)

# all src packages without vendor and generated code
PKGS = $(shell go list ./... | grep -v /vendor | grep -v /internal/server/restapi | grep -v /internal/server/grpcapi)

# Colors
GREEN_COLOR   = "\033[0;32m"
PURPLE_COLOR  = "\033[0;35m"
DEFAULT_COLOR = "\033[m"

all: clean fmt swagger build test

help:
	@echo 'Usage: make <TARGETS> ... <OPTIONS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    help               Show this help screen.'
	@echo '    clean              Remove binary.'
	@echo '    test               Run unit tests.'
	@echo '    lint               Run all linters including vet and gosec and others'
	@echo '    fmt                Run gofmt on package sources.'
	@echo '    build              Compile packages and dependencies.'
	@echo '    version            Print Go version.'
	@echo '    swagger            Generate swagger models and server'
	@echo '    swaggerdoc         Serve swagger doc'
	@echo ''
	@echo 'Targets run by default are: clean fmt lint test.'
	@echo ''

clean:
	@echo $(GREEN_COLOR)[clean]$(DEFAULT_COLOR)
	@$(GOCLEAN)
	@if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi

lint:
	@echo $(GREEN_COLOR)[lint]$(DEFAULT_COLOR)
	@$(GORUN) ./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint/main.go run \
	--no-config --disable=errcheck --enable=gosec --enable=prealloc ./...

test:
	@echo $(GREEN_COLOR)[test]$(DEFAULT_COLOR)
	@$(GOTEST) -race $(PKGS)

fmt:
	@echo $(GREEN_COLOR)[format]$(DEFAULT_COLOR)
	@$(GOFMT) $(PKGS)

build:
	@echo $(GREEN_COLOR)[build]$(DEFAULT_COLOR)
	$(GOBUILD) -o $(BINARY)

version:
	@echo $(GREEN_COLOR)[version]$(DEFAULT_COLOR)
	@$(GOCMD) version

swagger-clean:
	@echo $(GREEN_COLOR)[swagger cleanup]$(DEFAULT_COLOR)
	@rm -rf $(BASEPATH)/internal/server/models
	@rm -rf $(BASEPATH)/internal/server/restapi

swagger: swagger-clean swagger-build-binary
	@echo $(GREEN_COLOR)[swagger]$(DEFAULT_COLOR)
	@./bin/$(SWAGGER) generate server \
	   -f ./api/spec.yaml \
	   --exclude-main \
	   --flag-strategy=pflag \
	   --default-scheme=http \
	   --target=$(BASEPATH)/internal/server \
	   -q

swagger-build-binary:
ifeq ("$(wildcard ./bin/$(SWAGGER))","")
	@echo $(PURPLE_COLOR)[build swagger]$(DEFAULT_COLOR)
	@$(GOBUILD) -o ./bin/$(SWAGGER) ./vendor/github.com/go-swagger/go-swagger/cmd/swagger
endif

swaggerdoc: swagger
	@echo $(GREEN_COLOR)[doc]$(DEFAULT_COLOR)
	./bin/$(SWAGGER) serve --flavor=swagger $(BASEPATH)/api/spec.yaml

# start mySQL db in docker compose
mysql-up:
	@echo $(GREEN_COLOR)[start mysql]$(DEFAULT_COLOR)
	@$(DOCKER_COMPOSE) -f $(BASEPATH)/docker/storages/mysql.yml up -d

# stop mySQL db in docker compose
mysql-down:
	@echo $(GREEN_COLOR)[stop mysql]$(DEFAULT_COLOR)
	@$(DOCKER_COMPOSE) -f $(BASEPATH)/docker/storages/mysql.yml down