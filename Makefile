SWAGGER_VERSION = 2.2.10
API_VERSION ?= v0.1
MIGRATIONS_API_VERSION ?= v0.1
MIGRATIONS_REGISTRY_VERSION ?= v0.1
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
	docker build -t lendo/migrations-api:$(MIGRATIONS_API_VERSION) -f docker/migrations-api.Dockerfile .
	docker build -t lendo/migrations-registry:$(MIGRATIONS_REGISTRY_VERSION) -f docker/migrations-registry.Dockerfile .

.PHONY: build-image-api
build-image-api:
	docker build -t lendo/api:$(API_VERSION) -f docker/api.Dockerfile .

.PHONY: test-unit
test-unit:
	go test -v -count=1 ./api/...
	go test -v -count=1 ./registry/...

.PHONY: create-db-server
create-db-server:
	kubectl apply -f k8s/pg.yaml

.PHONY: create-db
create-db: create-db-server
	kubectl apply -f k8s/create_db.yaml

.PHONY: migrate-db
migrate-db: create-db
	kubectl apply -f k8s/migrations.yaml

.PHONY: create-api
create-api:
	kubectl apply -f k8s/api.yaml

.PHONY: docs
docs:
	mkdir -p api/docs

.PHONY: swagger
swagger: docs
	curl -s https://codeload.github.com/swagger-api/swagger-ui/tar.gz/v$(SWAGGER_VERSION) | tar xzv -C api/docs swagger-ui-$(SWAGGER_VERSION)/dist
	mv -f api/docs/swagger-ui-$(SWAGGER_VERSION)/dist/* api/docs/
	rm -rf api/docs/swagger-ui-$(SWAGGER_VERSION) api/docs/swagger-ui.js
	sed -i_ "s/swagger-ui\.js/swagger-ui\.min\.js/" api/docs/index.html
	sed -i_ "s/http:\/\/petstore\.swagger\.io\/v2\///" api/docs/index.html
	rm -f api/docs/*_

.PHONY: clean
clean:
	rm -rf api/docs
	find . -name '*.mock.go' -delete

.PHONY: generate
generate:
	go generate ./...
