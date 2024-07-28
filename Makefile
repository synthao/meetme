# Image and tag can be overridden via environment variables.
DOCKER_USERNAME ?= synthao
IMAGE_NAME ?= meetme
TAG ?= latest

# Name of the Docker image.
IMAGE := ${DOCKER_USERNAME}/${IMAGE_NAME}:${TAG}

APP_PORT?=8080

# Postgres
POSTGRES_CONTAINER_NAME?=meetme_pg
POSTGRES_VERSION=16.3
POSTGRES_USER?=root
POSTGRES_PASSWORD?=123456
POSTGRES_DB?=meetme
POSTGRES_PORT?=5432

.PHONY: app-docker-build app-docker-push dp pg-run pg-stop pg-rm migrate-up migrate-create-schema clean gen k-restart-pods

docker-build-push: app-docker-build app-docker-push

app-docker-build:
	@echo "Building Docker image ${IMAGE_NAME}"
	@docker build -t ${IMAGE} .

app-docker-push:
	@echo "Pushing Docker image ${IMAGE_NAME}"
	@docker push ${IMAGE}

app-docker-run:
	@echo "Running Docker image ${IMAGE_NAME}"
	@docker run -d --name $(IMAGE_NAME) -p $(APP_PORT):$(APP_PORT) --env-file .env $(IMAGE)

app-docker-stop:
	@echo "Stopping Docker image ${IMAGE_NAME}"
	@docker stop $(IMAGE_NAME)
	@docker rm $(IMAGE_NAME)
	@docker rmi $(IMAGE) -f

migrate-create-schema:
	@echo "Creating schema"
	@migrate create -ext sql -dir migrations -seq init_schema

migrate-up:
	@echo "Migrating up"
	@migrate -source file://migrations -database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable up

pg-run:
	@echo "Running Postgres container"
	@docker run --name $(POSTGRES_CONTAINER_NAME) \
		-e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_DB=$(POSTGRES_DB) \
		-p $(POSTGRES_PORT):5432 \
		-d postgres:$(POSTGRES_VERSION)

pg-start:
	@echo "Starting Postgres container"
	@docker start $(POSTGRES_CONTAINER_NAME)

pg-stop:
	@echo "Stopping Postgres container"
	@docker stop $(POSTGRES_CONTAINER_NAME)

pg-rm:
	@echo "Stopping Postgres container"
	@docker stop $(POSTGRES_CONTAINER_NAME)
	@echo "Removing Postgres image"
	@docker rm $(POSTGRES_CONTAINER_NAME)

clean: pg-stop pg-rm
	@echo "Removing all docker images"
	@docker stop $(IMAGE)
	@docker rmi $(IMAGE)
	@docker rmi $(POSTGRES_CONTAINER_NAME)

k-restart-pods:
	@kubectl apply -f deployments/
	@kubectl delete pods -l app=meetme

gen:
	@protoc -I proto proto/imgproc/imgproc.proto --go_out=./gen/go/ --go_opt=paths=source_relative --go-grpc_out=./gen/go/ --go-grpc_opt=paths=source_relative