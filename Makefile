COMPOSE_FILE := build/docker-compose.yaml
DEV_COMPOSE_FILE := build/docker-compose.dev.yaml
ENV_FILE := .env

run:
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up

build_app:
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up --build

dev:
	docker compose -f $(DEV_COMPOSE_FILE) --env-file $(ENV_FILE) up

stop:
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) down

delete:
	docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) down -v