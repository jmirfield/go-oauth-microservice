# main
run:
	go run ./cmd/server

test:
	go test -cover -v -race ./...

build:
	env GOOS=linux CGO_ENABLED=0 go build -o ./bin/oauth ./cmd/server/main.go

# docker compose
up: build
	docker-compose up --build -d

down:
	docker-compose down --remove-orphans


# docker support
clean:
	docker system prune -f