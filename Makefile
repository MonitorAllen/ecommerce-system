DB_URL=postgresql://root:secret@localhost:5432/ecom?sslmode=disable

create-db:
	docker exec -it postgres createdb --username=root --owner=root ecom

drop-db:
	docker exec -it postgres dropdb ecom

new-migration:
	 goose -s -dir "internal/db/migrations" create $(name) sql

migrate-up:
	goose up

migrate-up1:
	goose up-by-one

migrate-down:
	goose down

migrate-down1:
	goose down-by-one

sqlc:
	sqlc generate


.PHONY: new-migration sqlc create-db drop-db