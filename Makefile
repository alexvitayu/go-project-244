# Собирает бинарный файл в bin/gendiff
build:
	go build -o bin/gendiff ./cmd/gendiff

# Устанавливает собранный бинарник в GOBIN, чтобы его можно было запускать из любого места.
install: build
	go install ./cmd/gendiff

# Запуск линтера
lint:
	golangci-lint run

# Запуск тестов
test:
	go test -v ./...

