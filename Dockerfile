# Multi-stage build for skillKonnect

# Builder stage
FROM golang:1.24-bullseye AS builder
WORKDIR /src

# Download dependencies first (cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build with CGO enabled for sqlite3 (mattn/go-sqlite3)
ENV CGO_ENABLED=1
RUN mkdir -p /app/bin && \
	go build -ldflags "-s -w" -o /app/bin/skillkonnect ./cmd/app

# Final image: keep it minimal but compatible
FROM debian:bullseye-slim

# Install runtime dependency for sqlite and certs
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y --no-install-recommends libsqlite3-0 ca-certificates curl && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary and public assets from builder
COPY --from=builder /app/bin/skillkonnect /app/skillkonnect
COPY --from=builder /src/public /app/public

# Provide an empty .env to silence libraries expecting its presence
RUN touch /app/.env

# Data directory for sqlite DB (mount this as a volume)
RUN mkdir -p /data /app/public/uploads
VOLUME ["/data", "/app/public/uploads"]

# Default envs (override with docker-compose or CLI)
ENV HTTP_LISTEN_ADDR=":8080"
ENV DB_DRIVER=sqlite3
ENV DB_NAME=/data/app_db

EXPOSE 8080

ENTRYPOINT ["/app/skillkonnect"]
