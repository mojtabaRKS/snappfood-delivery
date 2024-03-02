COMPOSE_FILES=docker-compose.yml
COMPOSE_PROFILES=
COMPOSE_COMMAND=docker-compose

ifeq (, $(shell which $(COMPOSE_COMMAND)))
	COMPOSE_COMMAND=docker compose
	ifeq (, $(shell which $(COMPOSE_COMMAND)))
		$(error "No docker compose in path, consider installing docker on your machine.")
	endif
endif

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# If the first argument is "log"...
ifeq (log,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif


help:
	@echo "env"
	@echo "==> Create .env file"
	@echo ""
	@echo "up"
	@echo "==> Create and start containers"
	@echo ""
	@echo "build-up"
	@echo "==> Create and build all containers"
	@echo ""
	@echo "status"
	@echo "==> Show currently running containers"
	@echo ""
	@echo "destroy"
	@echo "==> Down all the containers, keeping their data"
	@echo ""
	@echo "purge"
	@echo "==> Down all the containers, removing their data"
	@echo ""
	@echo "psql-shell"
	@echo "==> Create an interactive shell for psql"
	@echo ""
env:
	@[ -e ./.env ] || cp -v ./.env.example ./.env

up:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) up -d

build-up:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) up -d --build --force-recreate

build-no-cache:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) build --no-cache

status:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) ps

down:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) down --remove-orphans

purge:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) down --remove-orphans --volumes

mysql-shell:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) exec -u 0 mysql mysql -hmysql -u$(MYSQL_USER) -D$(MYSQL_DATABASE) -p$(MYSQL_PASSWORD)

redis-shell:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) exec -u 0 redis redis-cli
	
.PHONY: log
log:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) logs -f $(RUN_ARGS)

.PHONY: log
shell:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) exec $(RUN_ARGS) bash