FROM golang:alpine

RUN apk update && apk add --no-cache git gcc sqlite-dev musl-dev

WORKDIR /app

COPY . .

RUN go install github.com/gocopper/cli/cmd/copper@latest
RUN --mount=type=cache,target=/root/.cache/go-build copper build

CMD ["/bin/sh", "-c", "/app/build/migrate.out --config ./config/prod.toml && /app/build/app.out  --config ./config/prod.toml"]