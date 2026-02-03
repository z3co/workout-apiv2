# Generate go from sql with sqlc
[group("sqlc")]
generate: 
	sqlc generate

# Test the package
[group("go")]
test: generate 
	go test ./...

# Run and test the package
[group("go")]
dev: generate test 
	go run main.go

# Start the local db
[group("db")]
up: 
	docker compose up -d

# Stop and remove the local db
[group("db")]
down: 
	docker compose down -v --rmi all

