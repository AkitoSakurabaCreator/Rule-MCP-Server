.PHONY: build test run clean

# ビルド
build:
	go build -o rule-mcp-server ./cmd/server

# テスト実行
test:
	go test -v ./...

# テストカバレッジ
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# サーバー起動
run:
	go run ./cmd/server

# カスタムポートでサーバー起動
run-port:
	@read -p "Enter port number: " port; \
	PORT=$$port go run ./cmd/server

# 本番環境でサーバー起動
run-prod:
	ENVIRONMENT=production LOG_LEVEL=warn go run ./cmd/server

# クリーンアップ
clean:
	rm -f rule-mcp-server
	rm -f coverage.out

# 依存関係の整理
deps:
	go mod tidy

# フォーマット
fmt:
	go fmt ./...

# リント
lint:
	golangci-lint run

# Docker コマンド
docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-clean:
	docker-compose down -v
	docker system prune -f

# ヘルプ
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  test         - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  run          - Run the server"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Tidy dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linter"
	@echo "  docker-build - Build Docker images"
	@echo "  docker-up    - Start Docker services"
	@echo "  docker-down  - Stop Docker services"
	@echo "  docker-logs  - Show Docker logs"
	@echo "  docker-clean - Clean Docker resources"
	@echo "  help         - Show this help"
