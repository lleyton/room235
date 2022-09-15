FROM golang:alpine

RUN apk update && apk add --no-cache git gcc sqlite-dev musl-dev

WORKDIR /app

COPY . .

RUN go install github.com/gocopper/cli/cmd/copper@latest
RUN copper build

CMD ["/bin/sh", "-c", "/app/build/migrate.out && /app/build/app.out"]