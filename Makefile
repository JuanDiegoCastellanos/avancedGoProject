postgres:
	docker run --name postgres_16_alpine -p 5433:5432 -e POSTGRES_PASSWORD=manolo221212 -e POSTGRES_USER=root -d postgres:16-alpine
createdb:
	docker exec -it postgres_16_alpine createdb --username=root --owner=root simple_posts
dropdb:
	docker exec -it postgres_16_alpine dropdb simple_posts
migrateup:
	migrate -path ./db/migration -database "postgresql://root:manolo221212@localhost:5433/simple_posts?sslmode=disable" -verbose up
migratedown:
	migrate -path ./db/migration -database "postgresql://root:manolo221212@localhost:5433/simple_posts?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test: 
	go test -v ./...

.PHONY: postgres createdb migratedown migrateup sqlc test
