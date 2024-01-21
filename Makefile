SHELL:=/bin/bash
STATICCHECK=$(shell which staticcheck)
BULID_COMMIT=$(shell git rev-parse HEAD)
BUILD_DATE=$(shell date +'%Y/%m/%d %H:%M:%S')
LDFLAGS=-ldflags "-X main.buildCommit=${BULID_COMMIT} -X 'main.buildDate=${BUILD_DATE}'"

.DEFAULT_GOAL := build

test:
	go test -cover ./...

vet:
	go vet ./...
	$(STATICCHECK) ./...

generate:
	go generate ./...
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/proto/*.proto

gvt: generate vet test

cover:
	go test -cover ./... -coverprofile cover.out
	go tool cover -html cover.out -o cover.html

build: gvt 
	go build  -C cmd/server .
	go build  -C cmd/client .

runserver:
	go run ./cmd/server

runclient:
	go run ./cmd/client

rundb:
	docker compose up -d
	
racetest:
	go test -v -race ./...
