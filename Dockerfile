FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go install github.com/gocopper/cli/cmd/copper@latest
RUN copper build

CMD "./build/migrate.out && ./build/app.out"