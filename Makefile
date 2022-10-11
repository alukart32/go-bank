.PHONY: init-db postgres-up postgres-stop create-db drop-db migrate-up migrate-down sqlc test help

pg_name = postgres14
pg_user = postgres
pg_user_pass = postgres
pg_image = postgres:14-alpine
pg_uri = localhost:5432
db_name = bank

help:
	@echo List of params:
	@echo    pg_name            - postgres docker container name (default: $(pg_name))
	@echo    pg_user            - postgres root user (default: $(pg_user))
	@echo    pg_user_pass       - postgres root user password (default: $(pg_user_pass))
	@echo    pg_image           - postgres docker image (default: $(pg_image))
	@echo    pg_uri             - postgres uri (default: $(pg_uri))
	@echo    db_name            - postgres main db (default: $(db_name))
	@echo List of commands:
	@echo   make init-db            - init a new db for service
	@echo   make postgres-up        - run postgres docker container
	@echo   make postgres-stop      - stop postgres docker container
	@echo   make create-db          - create db in postgres
	@echo   make drop-db            - drop db in postgres
	@echo   make migrate-up         - start db migration, src - ./migrations
	@echo   make migrate-down       - rollback db migration
	@echo   make sqlc               - generate go files from sql
	@echo   make test               - run all tests

init-db: create-db migrate-up

postgres-up:
	docker run --name $(pg_name) -e POSTGRES_USER=$(pg_user) -e POSTGRES_PASSWORD=$(pg_user_pass) -p 5432:5432 -d $(pg_image)

postgres-stop:
	docker stop $(pg_name)

create-db:
	docker exec -it $(pg_name) createdb --username=$(pg_user) --owner=$(pg_user) $(db_name)

drop-db:
	docker exec -it $(pg_name) dropdb --username=$(pg_user) $(db_name)

migrate-up:
	migrate -database "postgresql://$(pg_user):$(pg_user_pass)@$(pg_uri)/$(db_name)?sslmode=disable" -path migrations -verbose up

migrate-down:
	migrate -database "postgresql://$(pg_user):$(pg_user_pass)@$(pg_uri)/$(db_name)?sslmode=disable" -path migrations -verbose down

sqlc:
	docker run --rm -v ${CURDIR}:/src -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...
