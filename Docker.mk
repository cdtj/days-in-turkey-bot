export VERSION=$(shell cat VERSION)
export BOT_TOKEN=$(shell cat .token)

.PHONY: build stop start clean

build:
	docker build --platform linux/amd64 \
		--build-arg BOT_TOKEN="${BOT_TOKEN}" \
		-t turkeydays:${VERSION} .

stop:
	docker stop $(shell docker container ls -a -f label="tag=turkeydays" -f "status=running" -q)

start:
	docker run --restart=unless-stopped \
		-v /home/docker/shared/days-in-turkey-bot:/db \
		-d turkeydays:${VERSION}

clean:
	docker rm $(shell docker container ls -a -f label="tag=turkeydays" -f "status=exited" -q)
	docker image rm $(shell docker images -a -f label="tag=turkeydays" -q)