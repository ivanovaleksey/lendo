API_VERSION ?= v0.1
MIGRATIONS_VERSION ?= v0.1
COMPONENT_EXE = bin/$(COMPONENT)

ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.PHONY: build-all
build-all:
	make build-app COMPONENT=api
	make build-app COMPONENT=registry

.PHONY: build-app
build-app:
	go build -o $(COMPONENT_EXE) ./$(COMPONENT)/cmd/

.PHONY: build-image-migrations
build-image-migrations:
	docker build -t lendo/migrations:$(MIGRATIONS_VERSION) -f docker/migrations.Dockerfile .

.PHONY: build-image-api
build-image-api:
	docker build -t lendo/api:$(API_VERSION) -f docker/api.Dockerfile .

.PHONY: test-unit
test-unit:
	go test -v -count=1 ./api/...
	go test -v -count=1 ./registry/...

.PHONY: create-db
create-db:
	kubectl apply -f k8s/pg.yaml

.PHONY: migrate-db
migrate-db: create-db
	kubectl apply -f k8s/migrations.yaml

.PHONY: create-api
create-api:
	kubectl apply -f k8s/api.yaml
