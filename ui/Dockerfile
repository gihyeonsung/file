# ui
FROM node:24-alpine AS builder

WORKDIR /app

COPY . .

RUN npm ci && \
    npm run build

# runtime
FROM caddy:latest

COPY --from=builder /app/dist ./

CMD ["caddy", "file-server", "--listen", ":80", "--browse"]
