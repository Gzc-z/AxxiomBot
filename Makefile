# Container
APP   = axxiom

# Compose
DEV_COMPOSE := $(wildcard docker-compose.dev.yaml)
ifeq ($(DEV_COMPOSE),)
	COMPOSE = docker compose -f docker-compose.dev.yaml
else
    COMPOSE = docker compose -f docker-compose.yaml
endif

# config
MSG_ARGS = $(filter-out $@,$(MAKECMDGOALS))
MSG ?= $(MSG_ARGS)
PATHGO = cmd/bot
RUN = go run $(PATHGO)/main.go $(MSG_ARGS)

.PHONY: msg

# --- setup ---

run:
	@$(COMPOSE) up $(APP) -d --force-recreate
	@$(COMPOSE) logs $(APP) -f

up:
	@$(COMPOSE) up $(APP) -d

down: stop
stop:
	@$(COMPOSE) down $(APP)

logs:
	@$(COMPOSE) logs $(APP) -f


restart:
	@$(COMPOSE) restart $(APP)

# --- build ---

build:
	@go build -v -o axxiom ./cmd/bot


# --- cleanup ---

clean:
	@$(COMPOSE) down -v --remove-orphans
	@$(COMPOSE) rm -vsf

# --- tests ---

test:
	@go test -count=1 ./...

test-ping:
	@go test -v -count=1 -tags=integration -run '^TestPing' ./...

bash:
	@$(COMPOSE) exec $(APP) /bin/env bash

msg:
	@go run ./internal/scripts/msg.go
