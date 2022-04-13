PROJECT_NAME=babytl_backend

docker-build:
	docker-compose up --build

docker-start:
	docker-compose up -d

setup:
	go mod download

format:
	go fmt ./...

build:
	make setup && go build -o=./bin/$(PROJECT_NAME)

start:
	./bin/$(PROJECT_NAME)

dev:
	go run .

clear:
	go mod tidy