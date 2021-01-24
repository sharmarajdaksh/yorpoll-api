.PHONY: logs build run clean db_logs app_logs generate_swagger_json

PROJECT=yorpoll
TARGET=bin/$(PROJECT)
SRC=cmd/$(PROJECT)/main.go
BIN=$(TARGET)/$(PROJECT)
SWAGGER_YAML=api/rest/v1/swagger.yaml
SWAGGER_UI=swaggerui

export ENV=dev
export LOG_FILE=logs/server.log

# Values: mysql, mongo
export DATABASE_TYPE=mongo
export DATABASE_PORT=27017

export DATABASE_NAME=$(PROJECT)
export DATABASE_HOST=127.0.0.1
export DATABASE_USER=mongo
export DATABASE_PASSWORD=mongopassword

export SERVER_HOST=127.0.0.1
export SERVER_PORT=9011

export DATABASE_ROOT_PASSWORD=mongorootpassword
export MYSQL_IMAGE_TAG=8.0
export MONGO_IMAGE_TAG=4.4.3-bionic
export MONGO_DOCKER_SCRIPT_PATH=scripts/db/mongodb/docker
export MONGO_DOCKER_SCRIPT_NAME=user_init.sh

SERVER_COMPOSE=docker/docker-compose.server.yml

# Configure the database to use
ifeq ($(DATABASE_TYPE),mysql)
DB_COMPOSE=docker/docker-compose.mysql.yml
endif
ifeq ($(DATABASE_TYPE),mongo)
DB_COMPOSE=docker/docker-compose.mongodb.yml
endif

COMPOSE_BASE_COMMAND=docker-compose -f $(SERVER_COMPOSE) -f $(DB_COMPOSE)

build:
	$(COMPOSE_BASE_COMMAND) build

run: build
	$(COMPOSE_BASE_COMMAND) up -d
	
clean: 
	$(COMPOSE_BASE_COMMAND) down

logs:
	$(COMPOSE_BASE_COMMAND) logs -f

db_logs:
	$(COMPOSE_BASE_COMMAND) logs -f db

app_logs:
	$(COMPOSE_BASE_COMMAND) logs -f app

generate_swagger_json:
	python3 -c 'import sys, yaml, json; json.dump(yaml.load(sys.stdin), sys.stdout, indent=4)' < api/rest/v1/swagger.yaml | tee swaggerui/swagger.json api/rest/v1/swagger.json