PROJECT = yorpoll
TARGET = bin/$(PROJECT)
SRC = cmd/$(PROJECT)/main.go
BIN = $(TARGET)/$(PROJECT)

export ENV = dev
export LOG_FILE = log/server.log

# export DATABASE_TYPE = mysql
export DATABASE_TYPE = mongo
export DATABASE_NAME = $(PROJECT)
export DATABASE_PORT = 9012
export DATABASE_HOST = 127.0.0.1
export DATABASE_USER = mongo
export DATABASE_PASSWORD = mongopassword

export SERVER_HOST = 127.0.0.1
export SERVER_PORT = 9011

DATABASE_ROOT_PASSWORD = mongorootpassword
MYSQL_DOCKER_IMAGE = mysql:8.0
MONGO_DOCKER_IMAGE = mongo:4.4.3-bionic 
MYSQL_CONTAINER_NAME = $(PROJECT)_mysql
MONGO_CONTAINER_NAME = $(PROJECT)_mongo
MONGO_DOCKER_SCRIPT_PATH = scripts/db/mongodb/docker
MONGO_DOCKER_SCRIPT_NAME = user_init.sh

build: $(wilcard **/*.go)
	mkdir -p $(TARGET) && \
	go build -o $(BIN) $(SRC)


run: build db
	mkdir -p log/ && \
	chmod +x $(BIN) && \
	./$(BIN)

clean: rm-db 
	rm -rf $(TARGET)/*

# Configure the database to use
ifeq ($(DATABASE_TYPE),mysql)
db: db-mysql
rm-db: rm-db-mysql
endif
ifeq ($(DATABASE_TYPE),mongo)
db: db-mongo
rm-db: rm-db-mongo
endif

db-mysql:
	- docker run \
		--name $(MYSQL_CONTAINER_NAME) \
		-e MYSQL_ROOT_PASSWORD=$(DATABASE_ROOT_PASSWORD) \
		-e MYSQL_DATABASE=$(DATABASE_NAME) \
		-e MYSQL_USER=$(DATABASE_USER) \
		-e MYSQL_PASSWORD=$(DATABASE_PASSWORD) \
		-p $(DATABASE_PORT):3306 \
		-d $(MYSQL_DOCKER_IMAGE)

rm-db-mysql:
	- docker rm -f $(MYSQL_CONTAINER_NAME)

db-mongo:
	- docker run \
		--name $(MONGO_CONTAINER_NAME) \
		-e MONGO_INITDB_ROOT_USERNAME=mongoadmin \
		-e MONGO_INITDB_ROOT_PASSWORD=$(DATABASE_ROOT_PASSWORD) \
		-e MONGO_INITDB_DATABASE=$(DATABASE_NAME) \
		-e DATABASE_USER=$(DATABASE_USER) \
		-e DATABASE_PASSWORD=$(DATABASE_PASSWORD) \
		--mount type=bind,source="$(shell pwd)/$(MONGO_DOCKER_SCRIPT_PATH)"/$(MONGO_DOCKER_SCRIPT_NAME),target=/docker-entrypoint-initdb.d/$(MONGO_DOCKER_SCRIPT_NAME) \
		-p $(DATABASE_PORT):27017 \
		-d $(MONGO_DOCKER_IMAGE)


rm-db-mongo:
	- docker image rm -f $(MONGO_IMAGE_NAME) 