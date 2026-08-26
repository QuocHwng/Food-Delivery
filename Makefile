.PHONY: help infra-up infra-down up down migrate-up migrate-down test lint

# ─── Variables ────────────────────────────────────────────────────────────────
SERVICES := user-service order-service payment-service restaurant-service notification-service
MIGRATE  := migrate  # golang-migrate CLI

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ─── Infrastructure ───────────────────────────────────────────────────────────
infra-up: ## Start PostgreSQL + RabbitMQ (dev mode)
	docker compose -f docker-compose.dev.yml up -d
	@echo "✅ Infrastructure up. PG: localhost:5432 | RabbitMQ UI: http://localhost:15672"

infra-down: ## Stop infrastructure
	docker compose -f docker-compose.dev.yml down

# ─── Full stack ───────────────────────────────────────────────────────────────
up: ## Build and start all services
	docker compose up -d --build

down: ## Stop all services
	docker compose down

# ─── Migrations ───────────────────────────────────────────────────────────────
migrate-up: ## Run all migrations (up)
	@for svc in $(SERVICES); do \
		echo "▶ Migrating $$svc..."; \
		$(MIGRATE) -path infrastructure/migrations/$$svc -database "$(DATABASE_URL)" up; \
	done

migrate-down: ## Rollback all migrations (down 1 step)
	@for svc in $(SERVICES); do \
		echo "▶ Rolling back $$svc..."; \
		$(MIGRATE) -path infrastructure/migrations/$$svc -database "$(DATABASE_URL)" down 1; \
	done

# ─── Development ──────────────────────────────────────────────────────────────
run-gateway: ## Run API Gateway locally
	cd backend/gateway && go run cmd/main.go

run-user: ## Run User Service locally
	cd backend/services/user-service && go run cmd/main.go

run-order: ## Run Order Service locally
	cd backend/services/order-service && go run cmd/main.go

run-payment: ## Run Payment Service locally
	cd backend/services/payment-service && go run cmd/main.go

run-restaurant: ## Run Restaurant Service locally
	cd backend/services/restaurant-service && go run cmd/main.go

run-notification: ## Run Notification Service locally
	cd backend/services/notification-service && go run cmd/main.go

run-frontend: ## Run frontend dev server
	cd frontend && npm run dev

# ─── Testing ──────────────────────────────────────────────────────────────────
test: ## Run all Go tests
	@for svc in gateway $(addprefix services/, $(SERVICES)); do \
		echo "▶ Testing backend/$$svc..."; \
		cd backend/$$svc && go test ./... -v && cd -; \
	done

# ─── Linting ──────────────────────────────────────────────────────────────────
lint: ## Run golangci-lint on all services
	@for svc in gateway $(addprefix services/, $(SERVICES)); do \
		echo "▶ Linting backend/$$svc..."; \
		cd backend/$$svc && golangci-lint run ./... && cd -; \
	done

# ─── Utilities ────────────────────────────────────────────────────────────────
tidy: ## Run go mod tidy on all services
	@for svc in gateway $(addprefix services/, $(SERVICES)) shared; do \
		echo "▶ Tidying backend/$$svc..."; \
		cd backend/$$svc && go mod tidy && cd -; \
	done
