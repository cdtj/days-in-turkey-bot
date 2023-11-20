# golang stuff
FROM golang:latest

RUN mkdir /app
WORKDIR /app

COPY pkg/go.mod .
COPY pkg/go.sum .

RUN go mod download && go mod verify
COPY pkg .

RUN go test -v ./...

RUN GOOS=linux GOARCH=amd64 go build -o bot -a -v cmd/bot/v2/main.go

LABEL tag="turkeydays"
ARG BOT_TOKEN
ENV BOT_TOKEN=${BOT_TOKEN}

ENTRYPOINT ["./bot"]
