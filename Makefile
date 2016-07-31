
CONTAINER=pdevine/whale-test

SRC_DIR:=$(shell pwd)

SRC:=$(shell find ./ -name \*.go -print)

# Build the individual binaries based on source level dependencies
all: whale

whale: $(SRC)
	@echo "Building the whale"
	@(CGO_ENABLED=0 GOOS=linux godep go build -a -tags "netgo static_build" -installsuffix netgo -ldflags "-w" -o whale .)

docker: whale
	@echo "Building whale container"
	@docker build -t $(CONTAINER) -f Dockerfile ./
