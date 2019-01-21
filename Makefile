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

clean:
	@echo $(GREEN_COLOR)[clean]$(DEFAULT_COLOR)
	@$(GOCLEAN)
	@if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi
	@rm -rf ./bin

fmt:
	@echo $(GREEN_COLOR)[format]$(DEFAULT_COLOR)
	@$(GOFMT) $(PKGS)

swagger-clean:
	@echo $(GREEN_COLOR)[swagger cleanup]$(DEFAULT_COLOR)
	@rm -rf $(BASEPATH)/internal/server/models
	@rm -rf $(BASEPATH)/internal/server/restapi

swagger: swagger-clean swagger-build-binary
	@echo $(GREEN_COLOR)[swagger]$(DEFAULT_COLOR)
	@$(SWAGGER) generate server \
	   -f ./api/spec.yaml \
	   --exclude-main \
	   --flag-strategy=pflag \
	   --default-scheme=http \
	   --target=$(BASEPATH)/internal/server

swagger-build-binary:
ifeq ("$(wildcard ./bin/$(SWAGGER))","")
	@echo $(PURPLE_COLOR)[build swagger]$(DEFAULT_COLOR)
	@$(GOBUILD) -o ./bin/$(SWAGGER) ./vendor/github.com/go-swagger/go-swagger/cmd/swagger
endif

swaggerdoc: swagger
	@echo $(GREEN_COLOR)[doc]$(DEFAULT_COLOR)
	@$(SWAGGER) serve --flavor=swagger $(BASEPATH)/api/spec.yaml