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

dokku-deploy:
	$(if $(shell git config remote.dokku.url),, git remote add dokku $(DOKKU_URL))
	git push -f dokku $(shell git rev-parse HEAD):refs/heads/master