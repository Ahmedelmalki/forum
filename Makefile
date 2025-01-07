PORT = 8090
DB_PATH = data.db

build:
	@docker build -t forum-image .

run:
	@docker run --name forum-container -p $(PORT):$(PORT) forum-image

stop:
	@docker stop forum-container || true

clean:
	@rm -rf forum || true
	@docker rm -f forum-container || true
	@docker rmi -f forum-image || true

all: stop clean build run

run-app:
	@echo "Running application directly on host"
	@go run main.go
