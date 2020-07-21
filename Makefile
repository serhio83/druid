export WORKDIR := $(shell pwd)
export GO111MODULE=on
PREFIX?=serhio
APP?=druid
PORT?=9999
RELEASE?=v0.0.1
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
PROJECT?=github.com/serhio83/druid/pkg
GOOS?=linux
GOARCH?=amd64

clean:
	rm -f ${APP}

deps:
	go generate

build: clean
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
	-ldflags "-s -w -X ${PROJECT}/version.Release=${RELEASE} \
	-X ${PROJECT}/version.Commit=${COMMIT} -X ${PROJECT}/version.BuildTime=${BUILD_TIME}" \
	-o ${APP}

container: build
	docker build -t $(PREFIX)/$(APP):$(RELEASE) .

run: container
	docker run --name ${APP} -p ${PORT}:${PORT} --rm -e "PORT=${PORT}" -v $(WORKDIR)/.ssh:/root/.ssh $(PREFIX)/$(APP):$(RELEASE)

keys:
	ssh-keygen -t rsa -b 4096 -f .ssh/id_rsa -q -N ""
	ssh-copy-id -f -i .ssh/id_rsa root@gw.tp.fbs

push:
	docker push $(PREFIX)/$(APP):$(RELEASE)
