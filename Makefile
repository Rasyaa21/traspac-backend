.PHONY: \
  swag docs clean-docs install-swag \
  run dev build format tidy test \
  dev-up dev-down dev-down-v dev-logs dev-ps dev-rebuild dev-restart \
  prod-up prod-down prod-down-v prod-logs prod-ps prod-rebuild prod-restart

SWAG=swag
MAIN=cmd/server/main.go
DOCS=cmd/server/docs

DEV_COMPOSE=docker-compose.dev.yml
PROD_COMPOSE=docker-compose.prod.yml

# =========================
# Swagger
# =========================
swag:
	@echo "Generating swagger docs..."
	@rm -rf $(DOCS)
	@$(SWAG) init -g $(MAIN) -o $(DOCS) --parseDependency --parseInternal
	@echo "Swagger docs generated successfully!"

docs: swag

clean-docs:
	@echo "Cleaning swagger docs..."
	@rm -rf $(DOCS)
	@echo "Swagger docs cleaned!"

install-swag:
	@echo "Installing swag..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Swag installed successfully!"

# =========================
# Go App (Local non-docker)
# =========================
run:
	@echo "Starting server..."
	@go run $(MAIN)

dev: swag run

build:
	@echo "Building binary..."
	@mkdir -p bin
	@go build -o bin/server $(MAIN)

format:
	@echo "Formatting code..."
	@go fmt ./...

tidy:
	@echo "Tidying modules..."
	@go mod tidy

test:
	@echo "Running tests..."
	@go test ./... -v

# =========================
# Docker Compose - Development
# =========================
dev-up:
	@echo "Starting DEV stack..."
	@docker compose -f $(DEV_COMPOSE) up --build

dev-down:
	@echo "Stopping DEV stack..."
	@docker compose -f $(DEV_COMPOSE) down

dev-down-v:
	@echo "⚠️ DEV down -v (removing volumes / DB data)..."
	@docker compose -f $(DEV_COMPOSE) down -v

dev-logs:
	@docker compose -f $(DEV_COMPOSE) logs -f app

dev-ps:
	@docker compose -f $(DEV_COMPOSE) ps

dev-rebuild:
	@echo "Rebuilding DEV images (no cache)..."
	@docker compose -f $(DEV_COMPOSE) build --no-cache

dev-restart:
	@echo "Restarting DEV stack..."
	@docker compose -f $(DEV_COMPOSE) down
	@docker compose -f $(DEV_COMPOSE) up --build

# =========================
# Docker Compose - Production
# =========================
prod-up:
	@echo "Starting PROD stack (detached)..."
	@docker compose -f $(PROD_COMPOSE) up --build -d

prod-down:
	@echo "Stopping PROD stack..."
	@docker compose -f $(PROD_COMPOSE) down

prod-down-v:
	@echo "⚠️ PROD down -v will DELETE volumes / DB data."
	@docker compose -f $(PROD_COMPOSE) down -v

prod-logs:
	@docker compose -f $(PROD_COMPOSE) logs -f app

prod-ps:
	@docker compose -f $(PROD_COMPOSE) ps

prod-rebuild:
	@echo "Rebuilding PROD images (no cache)..."
	@docker compose -f $(PROD_COMPOSE) build --no-cache

prod-restart:
	@echo "Restarting PROD stack..."
	@docker compose -f $(PROD_COMPOSE) down
	@docker compose -f $(PROD_COMPOSE) up --build -d