APP=main
CONTAINERS=postgres postgres-pgadmin redis redis-commander
POSTGRES_CONTAINER=postgres
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USERNAME=admin

test:
	go test ./...

lint:
	gofmt -s -w .
	go vet ./...

clean:
	go clean -testcache ./...

migrate:
	go run cmd/migrate/main.go
