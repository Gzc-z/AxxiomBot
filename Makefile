# Container
APP   = axxiom
COMPOSE = docker compose -f docker-compose.yaml

# config
MSG_ARGS = $(filter-out $@,$(MAKECMDGOALS))
MSG ?= $(MSG_ARGS)
GOPATH = cmd/bot
RUN = go run $(GOPATH)/main.go $(MSG_ARGS)

.PHONY: msg

# --- setup ---

run:
	@$(COMPOSE) up $(APP) -d --force-recreate

up:
	@$(COMPOSE) up $(APP) -d

down: stop
stop:
	@$(COMPOSE) down $(APP)

logs:
	@$(COMPOSE) logs $(APP) -f


restart: build-app run
reset: clean build-app run

# --- build ---

build:
	@docker build -t $(APP) .

bash:
	@$(COMPOSE) exec $(APP) /bin/env bash

# --- cleanup ---

clean:
	@$(COMPOSE) down -v --remove-orphans
	@$(COMPOSE) rm -vsf

# --- tests ---

test:
	@go test -count=1 -tags=integration ./...

test-ping:
	@go test -v -count=1 -tags=integration -run '^TestPing' ./...

msg:
	@go test -v -count=1 -tags=integration -run '^TestPing' ./... -args "$(MSG)"
