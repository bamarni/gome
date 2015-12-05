PWD ?= $(shell pwd)

NS = bamarni
REPO = gome
VERSION ?= latest

NAME = bamarni-gome
INSTANCE ?= default
PORTS = -p 80:80
ENV = -e FORECAST_API_KEY -e VBB_API_KEY

.PHONY: build push pull run stop rm

build:
	mkdir -p build

	docker run --rm --name $(NAME)-$(INSTANCE)-node -v $(PWD)/web:/usr/src/app -w /usr/src/app \
		node:4-slim /bin/sh -c "npm i && npm run tsc && npm run uglifyjs"

	docker run --rm --name $(NAME)-$(INSTANCE)-golang -v $(PWD):/go/src/github.com/bamarni/gome -w /go/src/github.com/bamarni/gome \
		golang:1 /bin/sh -c "go get && CGO_ENABLED=0 go build -o ./build/gome -a -ldflags '-s' ./gome.go"

	docker create --name $(NAME)-$(INSTANCE)-jessie-curl buildpack-deps:jessie-curl &&\
		docker cp $(NAME)-$(INSTANCE)-jessie-curl:/etc/ssl/certs/ca-certificates.crt ./build/ca-certificates.crt &&\
		docker cp $(NAME)-$(INSTANCE)-jessie-curl:/usr/share/zoneinfo/Europe/Berlin ./build/zoneinfo-berlin &&\
		docker rm $(NAME)-$(INSTANCE)-jessie-curl

	docker build -t $(NS)/$(REPO):$(VERSION) .

push:
	docker push $(NS)/$(REPO):$(VERSION)

pull:
	docker pull $(NS)/$(REPO):$(VERSION)

run:
	docker run --rm --name $(NAME)-$(INSTANCE) $(PORTS) $(ENV) $(NS)/$(REPO):$(VERSION)

stop:
	docker stop $(NAME)-$(INSTANCE)

rm:
	docker rm $(NAME)-$(INSTANCE)

default: build
