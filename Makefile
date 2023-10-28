.PHONY: build build-prof build-run dc mocks test run lint

build:
	go build -o ./build/server ./cmd/server.go

build-prof: build
	go tool pprof —text ./bin/server

build-run: build
	./build/server

dc:
	docker-compose up  --remove-orphans --build

mocks:
	go generate ./...

test:
	go test -v -coverprofile cover.out ./... && go tool cover -html=cover.out

run:
	go run -race ./cmd/server.go

lint:
	golangci-lint run