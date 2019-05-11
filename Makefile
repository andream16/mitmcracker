build:
	go build -o mitmcracker cmd/main.go
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
