export VERSION=$(shell cat VERSION)
export BOT_TOKEN=$(shell cat .token)
export BOLDTB_PATH=./

.PHONY: build stop start clean

build:
	docker build --platform linux/amd64 \
		--build-arg BOT_TOKEN="${BOT_TOKEN}" \
		--build-arg BOLDTB_PATH="${BOLDTB_PATH}" \
		-t turkeydays:${VERSION} .

stop:
	docker stop $(shell docker container ls -a -f label="tag=turkeydays" -f "status=running" -q)

start:
	docker run --restart=unless-stopped \
		-d turkeydays:${VERSION}

clean:
	docker rm $(shell docker container ls -a -f label="tag=turkeydays" -f "status=exited" -q)
	docker image rm $(shell docker images -a -f label="tag=turkeydays" -q)