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
	docker build -t rrednoss/alertmanager-signl4:0.1.1 .
.PHONY:build

deploy: build
	docker push rrednoss/alertmanager-signl4:0.1.1
.PHONY: deploy

k8sDeploy: deploy
	helm upgrade -i -n alertmanager-signl4 alertmanager-signl4 ./chart/alertmanager-signl4
.PHONY: k8sDeploy

k8sPortForward: k8sDeploy
	kubectl port-forward svc/alertmanager-signl4 8888:80
.PHONY: k8sPortForward
