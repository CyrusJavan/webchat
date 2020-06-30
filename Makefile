.PHONY: all

all: build

build:
	docker build -t cjavan/caddy-go -f caddy/Dockerfile .
	docker push cjavan/caddy-go
	docker build -t cjavan/webchat-go -f ./app/Dockerfile .
	docker push cjavan/webchat-go