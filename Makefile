build:
	CGO_ENABLED=0 GO111MODULE=on go build -ldflags "-s -w" -o mitmcracker cmd/main.go

test:
	gotest ./... -cover -race -count=1
test-integration:
	docker-compose up -d redis 
	gotest ./... -cover -race -count=1 -tags=integration
	docker-compose stop redis
test-integration-cleanup:
	docker-compose down --rmi local
	docker-compose up -d redis 
	gotest ./... -cover -race -count=1 -tags=integration
	docker-compose stop redis
