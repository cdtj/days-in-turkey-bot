# golang stuff
FROM golang:latest as build-env

RUN mkdir /app
WORKDIR /app

COPY pkg/go.mod .
COPY pkg/go.sum .

RUN go mod download && go mod verify
COPY pkg .

RUN go test -v ./...

RUN GOOS=linux GOARCH=amd64 go build -o bot -a -v cmd/bot/v2/main.go

FROM gcr.io/distroless/base-debian12:latest AS build-release

# assets are ebedded
# COPY --from=build-env /app/assets /app/assets
COPY --from=build-env /app/bot /app/bot
WORKDIR /app

LABEL tag="turkeydays"
ARG BOT_TOKEN
ENV BOT_TOKEN=${BOT_TOKEN}
ARG BOLDTB_PATH
ENV BOLDTB_PATH=${BOLDTB_PATH}

ENTRYPOINT ["./bot"]