include .env

PACKAGES := $(shell go list ./...)
BASENAME := $(shell basename ${PWD})

.PHONY: help build run live connect db hasura redis create_network create_volume stop exec cleanI cleanC test
.DEFAULT_GOAL:= run

help: makefile
	@echo
	@echo " Choose a make command to run"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## build: Build the Docker image
build:
	@docker build -t $(APP_NAME) .

## run: Build the Docker image and run the container
run: cleanC create_network build
	@docker run -d \
		--env-file .env \
		--restart always \
		--name $(APP_NAME) \
		--network $(PROJECT_NAME) \
		-p $(APP_PORT):$(APP_PORT) \
		$(APP_NAME)

## live: Go build and running
live:
#	find . -type f \( -name '*.go' -o -name '*.gohtml' \) | entr -r sh -c 'make && docker logs --follow $(APP_NAME)'
	find . -type f \( -name '*.go' -o -name '*.gohtml' \) | entr -r sh -c 'go build -o /tmp/build ./cmd && /tmp/build'

## connect: 
connect:
ifeq ($(APP_ENV),local)
	$(MAKE) db
	$(MAKE) redis
	$(MAKE) hasura
endif

## db: 
db: create_volume
	docker run -d --name $(APP_NAME)-postgres --restart always \
		-p $(DB_PORT):$(DB_PORT) \
		-e TZ="Europe/Istanbul" \
		-e POSTGRES_DB=$(DB_NAME) \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASS) \
		-v $(APP_NAME):/var/lib/postgresql/data \
		--network $(PROJECT_NAME) \
		postgres:latest

hasura:
	docker run -d --name $(APP_NAME)-hasura --restart always \
		-p $(HASURA_PORT):$(HASURA_PORT) \
		-e HASURA_GRAPHQL_DATABASE_URL="postgres://$(DB_USER):$(DB_PASS)@$(APP_NAME)-postgres:$(DB_PORT)/$(DB_NAME)" \
		-e HASURA_GRAPHQL_ENABLED_LOG_TYPES="startup, http-log, webhook-log, websocket-log, query-log" \
		-e HASURA_GRAPHQL_ENABLE_CONSOLE="true" \
		-e HASURA_GRAPHQL_ADMIN_SECRET=$(HASURA_ADMIN_SECRET) \
		-e HASURA_GRAPHQL_JWT_SECRET='{"type":"HS256","key":"$(JWT_SECRET)"}' \
		-v $(APP_NAME):/hasura-migrations \
		--network $(PROJECT_NAME) \
		hasura/graphql-engine:v2.9.0

## redis: 
redis:
	docker run -d --name $(APP_NAME)-redis -p 6379:6379 --restart always --network $(PROJECT_NAME) redis:latest

## create_network: 
create_network:
	@if ! docker network inspect $(PROJECT_NAME) >/dev/null 2>&1; then \
		docker network create $(PROJECT_NAME); \
	fi

## create_volume: 
create_volume:
	@if ! docker volume inspect $(PROJECT_NAME) >/dev/null 2>&1; then \
		docker volume create $(PROJECT_NAME); \
	fi

## stop: Stop and remove the Docker container 
stop:
	@docker stop --time=600 $(APP_NAME)
	@docker rm $(APP_NAME)

## exec: Run the application inside the Docker container
exec:
	@docker exec -it $(APP_NAME) $(CMD)

## cleanI: Clean up the Docker image
cleanI:
	@docker rmi $(APP_NAME)
	@docker builder prune --filter="image=$(APP_NAME)"
	@docker rmi $(docker images -f "dangling=true" -q)

## cleanC: Clean up the Docker containers
cleanC:
	@CONTAINER_EXISTS=$$(docker ps -aq --filter name=$(APP_NAME)); \
	if [ -n "$$CONTAINER_EXISTS" ]; then \
		echo "Stopping and removing containers starting with $(APP_NAME)"; \
		CONTAINERS=$$(docker ps -aq --filter name=$(APP_NAME)); \
		for container in $$CONTAINERS; do \
			echo "Stopping and removing container $$container"; \
			docker stop $$container; \
			docker rm $$container; \
		done; \
	else \
		echo "No such containers starting with: $(APP_NAME)"; \
	fi

## test: Run all test
test: 
	@go test -v ./...