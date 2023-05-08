# basic parameters
NAME     := simple-line-ai-bot
VERSION  := v0.0.0
REVISION := $(shell git rev-parse --short HEAD)

# Go parameters
SRCS    := $(shell find . -type f -name '*.go')
DIST_DIRS := find * -type d -exec
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""
GOOS = "linux"
GOARCH = "amd64"
GOCMD = go
GOPRIVATE = github.com/swallowarc/*
GOBUILD = GOPRIVATE=$(GOPRIVATE) GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
GOVET = $(GOCMD) vet
GOGENERATE = $(GOCMD) generate
GOINSTALL = $(GOCMD) install

# build parameters
GRPC_PORT ?= 50051
DOCKER_CMD = docker
DOCKER_COMPOSE_CMD = docker-compose
DOCKER_BUILD = $(DOCKER_CMD) build
DOCKER_REGISTRY ?= fake-repository
DOCKER_USER ?= fake_user
DOCKER_PASS ?= fake_pass

# test parameters
MOCK_DIR=internal/tests/mocks/
REDIS_HOST_PORT?=localhost:6379

.PHONY: build setup-tools upgrade-grpc mock-clean mock-gen vet test docker/build
build:
	$(GOBUILD) -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/ -v ./...
setup/tools:
	$(GOINSTALL) github.com/golang/mock/mockgen@v1.6.0
setup/service:
ifeq ($(shell uname),Linux)
	$(DOCKER_COMPOSE_CMD) -f ./docker/docker-compose.yaml -f ./docker/docker-compose.override.yaml up -d
else
	$(DOCKER_COMPOSE_CMD) -f ./docker/docker-compose.yaml up -d
endif
mock-clean:
	rm -Rf ./$(MOCK_DIR)
mock-gen: mock-clean
	$(GOGENERATE) ./internal/usecases/...
	$(GOGENERATE) ./internal/interfaces/...
vet:
	$(GOVET) ./cmd/simple-line-ai-bot/...
test:
	$(GOTEST) -v ./...
docker/build:
	$(DOCKER_BUILD) -t $(DOCKER_REGISTRY) .
