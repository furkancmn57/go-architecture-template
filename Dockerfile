# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN swag init -g src/main.go -o src/docs --parseDependency --parseInternal --outputTypes go,json
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/bin/server ./src/main.go

FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -H -u 10001 appuser
COPY --from=builder /app/bin/server ./server

USER appuser
EXPOSE 7090

ENTRYPOINT ["./server"]
