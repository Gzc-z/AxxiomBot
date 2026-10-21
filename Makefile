# Container
APP   = axxiom
# COMPOSE = docker compose -f docker-compose.yaml
EXEC ?= /bin/env bash

# goConfig
GOPATH = cmd/bot
RUN = go run $(GOPATH)/main.go $(filter-out $@,$(MAKECMDGOALS))

.PHONY: test run ping

# --- setup ---

run:
	@$(COMPOSE) up $(APP) -d --force-recreate
	@$(COMPOSE) logs $(APP) -f

up:
	@$(COMPOSE) up $(APP) -d


stop:
	@$(COMPOSE) down $(APP)

build:
	@docker build -t $(APP) .


# ---

# run:
# 	@$(RUN)

test:
	@go test ./...

ping:
	@go test -v -count=1 ./internal

vet:
	@go vet ./...


# restart must build first
# hot-reload / live-reload does the job of restarting
restart: build-app run

reset: clean build-app run


# %:
# 	@:
