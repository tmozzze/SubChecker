.SILENT:

# env var
include .env
export

# Build binary
build:
	go mod download && go build -o ./bin/app ./cmd/app/main.go

# Run docker-compose
run:
	docker-compose up --build -d

# Down docker-compose
down:
	docker-compose down

# Down and clean docker-compose
down-and-clean:
	docker-compose down -v

# Create new migration (make create-migration NAME=name)
create-migration:
	migrate create -ext sql -dir database/migrations -seq $(NAME)

# Apply all migrations
migrate-up-all:
	migrate -path ./database/migrations -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" up

# Apply migration
migrate-up:
	migrate -path ./database/migrations -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" up 1

# Rollback last migration
migrate-down:
	migrate -path ./database/migrations -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" down 1

# Show migration status
migrate-status:
	@echo "migrate status:"
	migrate -path ./database/migrations -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" version

# Rollback all migrations
migrate-reset:
	migrate -path ./database/migrations -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" down

# Swagger docs gen
swagger-gen:
	swag init -g cmd/app/main.go -o docs

# First start Service
start-service: run