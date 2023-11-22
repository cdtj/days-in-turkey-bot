export VERSION=$(shell cat VERSION)
export BOT_TOKEN=$(shell cat .token)
export BOLDTB_PATH=./

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