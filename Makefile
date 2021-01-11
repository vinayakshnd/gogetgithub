BASE_DIR := $(shell git rev-parse --show-toplevel)
DOCKER_IMAGE := vinayakinfrac/gogetgithub:latest

BIN :=$(BASE_DIR)/bin

ifndef PKGS
# shell does not honor export command above, so we need to explicitly pass GOFLAGS here
PKGS := $(shell GOFLAGS=-mod=vendor go list ./... 2>&1)
endif

deps:
	go get -d -v $(PKGS)

build:
	mkdir -p $(BIN)
	go build .
	cp gogetgithub $(BIN)
	chmod -R 755 $(BIN)/*

test:
	go test ./httphandlers
	go test ./utils

container:
	docker build --tag $(DOCKER_IMAGE) -f Dockerfile .

deploy: container
	sudo docker push $(DOCKER_IMAGE)
