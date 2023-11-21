TERMSITE_VER=$(shell cat VERSION)
BOT_TOKEN=$(shell cat .token)

.PHONY: build run debug docker

build:
	cd ./pkg && \
	go build -o bot cmd/bot/v2/main.go

run:
	rm -f ./pkg/bot
	$(MAKE) build
	cd ./pkg && ./bot

debug:
	cd ./pkg && \
	go run ./cmd/bot/v2/main.go

docker:
	docker build --platform linux/amd64 \
		--build-arg BOT_TOKEN="${BOT_TOKEN}" \
		-t turkeydays:${VERSION} . && \
	docker run --restart=unless-stopped \
		-v /home/docker/shared/days-in-turkey-bot:/db \
		-d turkeydays:${VERSION}