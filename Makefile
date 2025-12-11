.PHONY: docs clean-docs run dev install-swag format tidy build test

SWAG=swag
MAIN=cmd/server/main.go
DOCS=cmd/server/docs

docs:
	@echo "Generating swagger docs..."
	@rm -rf $(DOCS)
	$(SWAG) init -g $(MAIN) -o $(DOCS) --parseInternal --parseDependency
	@echo "Swagger docs generated successfully!"

clean-docs:
	@echo "Cleaning swagger docs..."
	@rm -rf $(DOCS)
	@echo "Swagger docs cleaned!"

run:
	@echo "Starting server..."
	go run $(MAIN)

dev: docs run

install-swag:
	@echo "Installing swag..."
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Swag installed successfully!"

format:
	@echo "Formatting code..."
	go fmt ./...

tidy:
	@echo "Tidying modules..."
	go mod tidy

build:
	@echo "Building binary..."
	go build -o bin/server $(MAIN)

test:
	@echo "Running tests..."
	go test ./... -v