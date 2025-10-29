.SILENT:

# env var
include .env
export

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
	migrate -path ./database/migrations -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" version

# Rollback all migrations
migrate-reset:
	migrate -path ./database/migrations -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" down