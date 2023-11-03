NGROK_URL=https://3a22-188-169-160-237.ngrok.io

export BOT_TOKEN=$(shell cat .token)
export BOT_WEBHOOK=${NGROK_URL}/telegram-webhook

.PHONY: build run debug

build:
	cd ./src && \
	go build ${LDFLAGS} -o . 

run:
	rm -f ./src/termsite
	$(MAKE) -f Golang.mk build && \
	cd src && ./termsite

debug:
	cd ./pkg && \
	go run ./cmd/bot/main.go