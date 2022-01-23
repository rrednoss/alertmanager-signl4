.DEFAULT_GOAL := build

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

vet: fmt
	go vet ./...
	shadow ./...
.PHONY:vet

run:
	go run ./...
.PHONY:run

test:
	go test ./...
.PHONY:test

build: vet
	docker build -t rrednoss/alertmanager-signl4:0.1.0 .
.PHONY:build

deploy: build
	docker push rrednoss/alertmanager-signl4:0.1.0 .
.PHONY: deploy
