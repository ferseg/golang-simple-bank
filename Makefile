postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=sa123 -d postgres:15-alpine

create-db:
	docker exec -it postgres createdb --username=postgres --owner=postgres simple_bank

migrate-up:
	migrate -path db/migration -database "postgresql://postgres:sa123@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgresql://postgres:sa123@localhost:5432/simple_bank?sslmode=disable" -verbose down

delete-db:
	docker exec -it postgres dropdb simple_bank

.PHONY: postgres create-db drop-db migrate-up migrate-down
