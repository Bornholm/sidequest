SHELL := /bin/bash
DOKKU_URL := dokku@dev.lookingfora.name:sidequest

build: build-server build-client

build-server:
	go build -v -o bin/sidequest ./cmd/sidequest

build-client: node_modules
	npm run build

node_modules:
	npm ci

docker-image:
	docker build -t sidequest:latest .

docker-run:
	docker run \
		-it --rm \
		-p 8090:8090 \
		-v sidequest_data:/app/data \
		sidequest:latest

run-server:
	$(MAKE) run-with-env RUN_CMD="go run ./cmd/sidequest serve"

run-client:
	$(MAKE) run-with-env RUN_CMD="npm run watch"

run-with-env: .env
	( set -o allexport && source .env && $(RUN_CMD) )

.env:
	cp .env.dist .env

dokku-deploy:
	$(if $(shell git config remote.dokku.url),, git remote add dokku $(DOKKU_URL))
	git push -f dokku $(shell git rev-parse HEAD):refs/heads/master