# main
run:
	go run ./...

test:
	go test -cover -v -race ./...


# docker compose
up:
	docker-compose up --build -d

down:
	docker-compose down --remove-orphans


# docker support
clean:
	docker system prune -f