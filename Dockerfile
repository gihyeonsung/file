# file
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY . .

RUN apk add --no-cache \
    gcc \
    musl-dev \
    git \
    make \
    bash && \
    go build -o file .

# runtime
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/file ./file

RUN mkdir db data

CMD ["./file"]
