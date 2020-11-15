default: build

build:
	@'rm' ./resources/*.cmo ./resources/*.cmi ./resources/*.cmx ./resources/*.o ./resources/encrypt ./resources/decrypt 2>/dev/null || true
	ocamlopt -o ./resources/encrypt ./resources/encrypt.ml
	ocamlopt -o ./resources/decrypt ./resources/decrypt.ml
	CGO_ENABLED=0 GO111MODULE=on go build -ldflags "-s -w" -o mitmcracker cmd/main.go

test:
	go test ./... -count=10 -race