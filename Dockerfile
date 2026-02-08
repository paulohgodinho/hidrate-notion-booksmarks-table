# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY src/go.mod src/go.sum ./
RUN go mod download

# Copy source code
COPY src/ .

# Build the processor
RUN CGO_ENABLED=0 GOOS=linux go build -o /processor ./main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /processor .

# Copy .env file (will be overridden by docker-compose env_file)
COPY .env .env

CMD ["./processor"]
