APP_NAME=go-base-project
MAIN_FILE=./cmd/api/main.go
BINARY=./bin/$(APP_NAME)

.PHONY: help run build test clean migrate-up migrate-down tidy docker-up

help:
	@echo "Perintah yang tersedia:"
	@echo "  make run"
	@echo "  make build"
	@echo "  make test"
	@echo "  make test-cover"
	@echo "  make tidy"
	@echo "  make migrate-up"
	@echo "  make docker-up"
	@echo "  make clean"

run:
	@echo "Menjalankan aplikasi..."
	air

run-simple:
	go run $(MAIN_FILE)

build:
	@echo "Building..."
	@mkdir -p bin
	go build -o $(BINARY) $(MAIN_FILE)
	@echo "Binary tersimpan di $(BINARY)"

test:
	@echo "Menjalankan test..."
	go test ./... -v

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Laporan coverage: coverage.html"

tidy:
	go mod tidy

migrate-up:
	@echo "Menjalankan migrasi..."
	migrate -path ./migrations -database "mysql://root:@tcp(localhost:3306)/go_base_project" up

migrate-down:
	migrate -path ./migrations -database "mysql://root:@tcp(localhost:3306)/go_base_project" down 1

docker-up:
	@echo "Menjalankan MySQL..."
	docker run --name mysql-dev \
		-e MYSQL_ROOT_PASSWORD=root \
		-e MYSQL_DATABASE=go_base_project \
		-p 3306:3306 \
		-d mysql:8

docker-down:
	docker stop mysql-dev && docker rm mysql-dev

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html