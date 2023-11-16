postgres:
	docker run --name rssagg -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -p 5432:5432 -d postgres:16-alpine

createdb:
	docker exec -it rssagg createdb --username=root --owner=root rssagg

dropdb:
	docker exec -it rssagg dropdb rssagg

sqlc:
	sqlc generate

migrate-up:
	migrate -path sql/schema -database "postgresql://root:root@localhost:5432/rssagg?sslmode=disable" up

migrate-down:
	migrate -path sql/schema -database "postgresql://root:root@localhost:5432/rssagg?sslmode=disable" down

run:
	go build && ./RSSAggregator


.PHONY: postgres createdb dropdb sqlc migrate-up migrate-down run