include .env # make your conf from .env.example
export WORKDIR := $(shell pwd)
export GO111MODULE=on

PREFIX?=serhio
APP?=druid
RELEASE?=v0.0.1
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
PROJECT?=github.com/serhio83/druid/pkg
GOOS?=linux
GOARCH?=amd64

clean:
	rm -f ${APP}

build: clean
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
	-ldflags "-s -w -X ${PROJECT}/version.Release=${RELEASE} \
	-X ${PROJECT}/version.Commit=${COMMIT} -X ${PROJECT}/version.BuildTime=${BUILD_TIME}" \
	-o ${APP}

container: build
	docker build -t $(PREFIX)/$(APP):$(RELEASE) .

run: container
	docker run --name ${APP} -p ${DRUID_LISTEN_PORT}:${DRUID_LISTEN_PORT} --rm \
	-e "PORT=${DRUID_LISTEN_PORT}" \
	-e "DRUID_REGISTRY_HOST=${DRUID_REGISTRY_HOST}" \
	-e "DRUID_REGISTRY_PORT=${DRUID_REGISTRY_PORT}" \
	-e "DRUID_REGISTRY_USER=${DRUID_REGISTRY_USER}" \
	-e "DRUID_REGISTRY_PASSWORD=${DRUID_REGISTRY_PASSWORD}" \
	-e "DRUID_LISTEN_PORT=${DRUID_LISTEN_PORT}" \
	-e "DRUID_DATA_PATH=${DRUID_DATA_PATH}" \
	-v "${WORKDIR}/storage:/opt/druid" \
	$(PREFIX)/$(APP):$(RELEASE)

push:
	docker push $(PREFIX)/$(APP):$(RELEASE)
