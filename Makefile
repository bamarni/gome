PWD ?= $(shell pwd)

NS = bamarni
REPO = gome
VERSION ?= latest

NAME = bamarni-gome
INSTANCE ?= default
PORTS = -p 80:8080
VOLUMES = -v $(PWD)/web:/var/www
ENV = -e FORECAST_API_KEY

.PHONY: build push shell run start stop rm release

build:
	docker run --rm --name $(NAME)-node-$(INSTANCE) -v $(PWD)/web:/usr/src/app -w /usr/src/app node:4 /bin/sh -c "npm i && ./node_modules/.bin/tsc"
	docker build -t $(NS)/$(REPO):$(VERSION) .

push:
	docker push $(NS)/$(REPO):$(VERSION)

shell:
	docker run --rm --name $(NAME)-$(INSTANCE) -i -t $(PORTS) $(VOLUMES) $(ENV) $(NS)/$(REPO):$(VERSION) /bin/bash

run:
	docker run --rm --name $(NAME)-$(INSTANCE) $(PORTS) $(VOLUMES) $(ENV) $(NS)/$(REPO):$(VERSION)

start:
	docker run -d --name $(NAME)-$(INSTANCE) $(PORTS) $(VOLUMES) $(ENV) $(NS)/$(REPO):$(VERSION)

stop:
	docker stop $(NAME)-$(INSTANCE)

rm:
	docker rm $(NAME)-$(INSTANCE)

release: build
	make push -e VERSION=$(VERSION)

default: build