default: all

all:
	make test
	make build

build:
	@'rm' ./resources/*.cmo ./resources/*.cmi ./resources/*.cmx ./resources/*.o ./resources/encrypt ./resources/decrypt 2>/dev/null || true
	ocamlopt -o ./resources/encrypt ./resources/encrypt.ml
	ocamlopt -o ./resources/decrypt ./resources/decrypt.ml
	CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=$(shell git describe --tags --always --dirty)" -trimpath -o mitmcracker cmd/main.go

test:
	go test ./... -count=1 -race -coverprofile=coverage.out

coverage: test
	go tool cover -html=coverage.out -o coverage.html

lint:
	go vet ./...
	go fmt ./...

clean:
	@'rm' -f mitmcracker coverage.out coverage.html 2>/dev/null || true
	@'rm' ./resources/*.cmo ./resources/*.cmi ./resources/*.cmx ./resources/*.o ./resources/encrypt ./resources/decrypt 2>/dev/null || true

.PHONY: all build test coverage lint clean