# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main ./cmd/server/main.go

# Run stage
FROM alpine:latest
WORKDIR /app

# Copy built binary from builder
COPY --from=builder /app/main .
# Assuming .env is used in development but we can also use env vars in docker-compose.
# We will rely on environment variables passed by docker-compose in production,
# but we can copy .env just in case it's required for defaults.
COPY .env .env

# Expose port
EXPOSE 8080

CMD ["./main"]
