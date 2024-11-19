# Stage 1: Build the application
FROM golang:1.22 AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files for dependency resolution
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o url-shortener

# Stage 2: Create a minimal image for running the application
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/url-shortener .
COPY --from=builder /app/static ./static
COPY --from=builder /app/.env .env

# Expose application port
EXPOSE 8080

# Command to run the application
CMD ["./url-shortener"]
