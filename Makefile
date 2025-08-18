## Команда сборки исполняемого файла
build:
	go build -o bin/gendiff ./cmd/gendiff

install: build
	go install ./cmd/gendiff
