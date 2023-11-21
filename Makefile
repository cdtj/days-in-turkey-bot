NGROK_URL=https://2562-188-169-160-237.ngrok.io

export VERSION="0.0.1"
export BOT_TOKEN=$(shell cat .token)
export BOT_WEBHOOK=${NGROK_URL}/telegram-webhook

.PHONY: build run debug

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
	docker run --restart=unless-stopped turkeydays:${VERSION}