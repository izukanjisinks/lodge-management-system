# syntax=docker/dockerfile:1

# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Tidy modules and build
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o /app/hr-system ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests (email, etc.)
RUN apk --no-cache add ca-certificates

# Copy binary from builder stage
COPY --from=builder /app/hr-system /app/hr-system

# Expose port (actual port is set by SERVER_PORT env var at runtime)
EXPOSE 8081

# Run the application
CMD ["/app/hr-system"]
