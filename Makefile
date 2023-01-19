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

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mockgen:
	mockgen -package mockdb -destination db/mock/store.go github.com/ferseg/golang-simple-bank/db/sqlc Store

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		proto/*.proto

evans:
	evans --port 8080 --host localhost -r repl

.PHONY: postgres create-db drop-db migrate-up migrate-down sqlc mockgen proto evans
