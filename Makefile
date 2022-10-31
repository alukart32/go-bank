network = input_traffic
db_name = postgres14
db_user = postgres
db_user_pass = postgres
db_image = postgres:14-alpine
db_uri = localhost:5432
db_name = bank

.PHONY: help
help:
	@echo List of params:
	@echo    db_name            - postgres docker container name (default: $(db_name))
	@echo    db_user            - postgres root user (default: $(db_user))
	@echo    db_user_pass       - postgres root user password (default: $(db_user_pass))
	@echo    db_image           - postgres docker image (default: $(db_image))
	@echo    db_uri             - postgres uri (default: $(db_uri))
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

.PHONY: create-net
create-net:
	docker network create $(network)

.PHONY: compose-down
compose-down:
	docker compose -p account_app -f ./deployments/docker-compose.yml down

.PHONY: up-stage
up-stage:
	docker compose -p account_app -f ./deployments/docker-compose.yml --profile stage up -d

.PHONY: up-test
up-test:
	docker compose -p account_app -f ./deployments/docker-compose.yml --profile test up -d

.PHONY: migrate-up
migrate-up:
	migrate -database "postgresql://$(db_user):$(db_user_pass)@$(db_uri)/$(db_name)?sslmode=disable" -path migrations -verbose up

.PHONY: migrate-down
migrate-down:
	migrate -database "postgresql://$(db_user):$(db_user_pass)@$(db_uri)/$(db_name)?sslmode=disable" -path migrations -verbose down

.PHONY: sqlc
sqlc:
	docker run --rm -v ${CURDIR}:/src -w /src kjconroy/sqlc generate

.PHONY: test
test:
	go test -v -cover ./...
