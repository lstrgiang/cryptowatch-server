include .env
export 

docker/up:
	docker-compose up -d 

docker/down:
	docker-compose down

db/up:
	migrate -path ./db/migrations -database postgres://${DATABASE_USERNAME}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable up

db/down:
	migrate -path ./db/migrations -database postgres://${DATABASE_USERNAME}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable down
db/drop:
	migrate -path ./db/migrations -database postgres://${DATABASE_USERNAME}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable drop -f 

run:
	go run cmd/server/server.go
