BIN := "./bin/calendar"
BIN_SCHEDULER := "./bin/calendar_scheduler"
BIN_SENDER := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build: build_calendar build_scheduler build_sender

build_calendar:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run_calendar: build_calendar
	$(BIN) -config ./configs/config.toml

build_scheduler:
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler

run_scheduler: build_scheduler
	$(BIN_SCHEDULER) -config ./configs/config_scheduler.toml

build_sender:
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run_sender: build_sender
	$(BIN_SENDER) -config ./configs/config_sender.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

build-img_scheduler:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/DockerfileScheduler .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.45.2

lint: install-lint-deps
	golangci-lint run ./...

generate:
	protoc \
		--proto_path=api/ \
		--go_out=api/pb \
		--go-grpc_out=api/pb \
		--grpc-gateway_out=api/pb --grpc-gateway_opt paths=source_relative \
		api/*.proto

up:
	 docker-compose -f deployments/docker-compose.yaml -f deployments/docker-compose.dev.yaml up -d

down:
	 docker-compose -f deployments/docker-compose.yaml -f deployments/docker-compose.dev.yaml down

integration-tests:
	set -e; \
	docker-compose -f deployments/docker-compose.yaml -f deployments/docker-compose.test.yaml up --build -d; \
	test_status_code=0 ;\
	docker-compose -f deployments/docker-compose.yaml -f deployments/docker-compose.test.yaml run integration-tests go test -v -mod=readonly -tags integration ./tests/integration || test_status_code=$$? ;\
	docker-compose -f deployments/docker-compose.yaml -f deployments/docker-compose.test.yaml down ;\
	exit $$test_status_code ;

test-cleanup:
	docker-compose -f deployments/docker-compose.yaml -f deployments/docker-compose.test.yaml down \
	--rmi local \
			--volumes \
			--remove-orphans \
			--timeout 60; \
  	docker-compose rm -f

.PHONY: build run build-img run-img version test lint generate up down integration-tests


