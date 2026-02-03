generate:
	sqlc generate

db-up:
	docker compose up -d

db-down:
	docker compose down -v --rmi all

test: generate
	go test ./...

dev: generate test
	go run main.go
