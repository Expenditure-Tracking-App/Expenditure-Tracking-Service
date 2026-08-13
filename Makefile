APP_NAME=text-classifier-for-expenditure-tracking-bot
VENV=.venv
PYTHON=$(VENV)/bin/python
PIP=$(VENV)/bin/pip
UVICORN=$(VENV)/bin/uvicorn

# === Python Text Classifier Service ===
# Install dependencies into venv
install:
	python3 -m venv $(VENV)
	$(PIP) install --upgrade pip
	$(PIP) install -r requirements.txt

# Run the FastAPI service
run:
	$(UVICORN) text_classifier.app:app --reload --host 0.0.0.0 --port 8000

# Freeze current packages
freeze:
	$(PIP) freeze > requirements.txt

# Clean up the virtual environment
clean:
	rm -rf $(VENV)

# Run service + auto-install if needed
start: install run

# === Go Application ===
# Run the HTTP server (port 8081)
run-server:
	go run cmd/server/main.go

# Run the Telegram bot
run-bot:
	go run cmd/bot/main.go

# Build the project
build:
	go build -v ./...

# Run tests
test:
	go test ./...

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean-go:
	go clean ./...
