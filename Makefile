NGROK_URL=https://2562-188-169-160-237.ngrok.io

export BOT_TOKEN=$(shell cat .token)
export BOT_WEBHOOK=${NGROK_URL}/telegram-webhook

.PHONY: build run debug

build:
	cd ./pkg && \
	go build -o httpsrvr cmd/http-server/main.go

run:
	rm -f ./pkg/httpsrvr
	$(MAKE) build
	cd ./pkg && ./httpsrvr

debug:
	cd ./pkg && \
	go run ./cmd/bot/main.go

debugv2:
	cd ./pkg && \
	go run ./cmd/botv2/main.go