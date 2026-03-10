-include .env
export

APP_NAME=hr-system
MAIN=./cmd/api/main.go
MIGRATE_DIR=./migrations
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: build run dev tidy test migrate-up migrate-down migrate-create migrate-force \
        docker-build docker-up docker-up-d docker-stop docker-down docker-logs \
        docker-rebuild docker-ps docker-shell docker-restart docker-clean

# -----------------------------
# build: Compiles the Go source code into a binary
# -----------------------------
build:
	go build -o bin/$(APP_NAME) $(MAIN)

# -----------------------------
# run: Runs the Go project directly (without building a separate binary)
# -----------------------------
run:
	go run $(MAIN)

# -----------------------------
# dev: Alias for run
# -----------------------------
dev:
	go run $(MAIN)

# -----------------------------
# tidy: Tidy go modules
# -----------------------------
tidy:
	go mod tidy

# -----------------------------
# test: Run all tests
# -----------------------------
test:
	go test ./...

# -----------------------------
# Database migrations
# -----------------------------
migrate-up:
	migrate -path $(MIGRATE_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATE_DIR) -database "$(DB_URL)" down

migrate-create:
	migrate create -ext sql -dir $(MIGRATE_DIR) -seq $(name)

migrate-force:
	migrate -path $(MIGRATE_DIR) -database "$(DB_URL)" force $(version)

# -----------------------------
# Docker commands
# -----------------------------

docker-build:
	docker build -t $(APP_NAME):latest .

docker-up:
	docker-compose up

docker-up-d:
	docker-compose up -d

docker-stop:
	docker-compose stop

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-rebuild:
	docker-compose up --build

docker-ps:
	docker-compose ps

docker-shell:
	docker-compose exec api sh

docker-restart: docker-build docker-down docker-up-d

docker-clean:
	docker-compose down -v
	docker rmi $(APP_NAME):latest
