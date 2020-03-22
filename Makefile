all: build run

run:
	./divulge

build:
	go build -o divulge cmd/server/main.go

test:
	go test -v ./...

test_all:
	go test -v -tags integration ./...

start_pg:
	docker run -p 5432:5432 -e POSTGRES_PASSWORD=password postgres

migrate_up:
	go run cmd/migrate/migrate.go

migrate_down:
	go run cmd/migrate/migrate.go --down

cover:
	go test -coverprofile coverage.out ./...

cover_view: cover
	go tool cover -html coverage.out
