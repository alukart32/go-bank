.PHONY: init-db postgres-up postgres-stop create-db drop-db migrate-up migrate-down sqlc test help

network = input_traffic
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
	@echo    network            - project external network (default: $(network))
	@echo List of commands:
	@echo   make create-net         - create new docker net (default: $(network))
	@echo   make compose-down       - docker compose down
	@echo   make up-stage           - docker compose up stage profile
	@echo   make up-test            - docker compose up test profile
	@echo   make migrate-up         - start db migration, src - ./migrations
	@echo   make migrate-down       - rollback db migration
	@echo   make sqlc               - generate go files from sql
	@echo   make prepare-test       - prepare before test: up-test and migrate-up
	@echo   make test               - run all tests

create-net:
	docker network create $(network)

compose-down:
	docker compose -p account_app -f ./deployments/docker-compose.yml down

up-stage:
	docker compose -p account_app -f ./deployments/docker-compose.yml --profile stage up -d

prepare-test: up-test migrate-up

up-test:
	docker compose -p account_app -f ./deployments/docker-compose.yml --profile test up -d

migrate-up:
	migrate -database "postgresql://$(pg_user):$(pg_user_pass)@$(pg_uri)/$(db_name)?sslmode=disable" -path migrations -verbose up

migrate-down:
	migrate -database "postgresql://$(pg_user):$(pg_user_pass)@$(pg_uri)/$(db_name)?sslmode=disable" -path migrations -verbose down

sqlc:
	docker run --rm -v ${CURDIR}:/src -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...
