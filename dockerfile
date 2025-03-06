# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files from the microservice directory
COPY ../microservices/lgm8-notification-service/go.mod ../microservices/lgm8-notification-service/go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the microservice directory to the container
COPY ../microservices/lgm8-notification-service/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o lgm8-notification-service ./cmd/main.go

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Copy the built binary from the builder stage
COPY --from=builder /app/lgm8-notification-service /app/lgm8-notification-service

# Copy all config files
COPY ../microservices/lgm8-notification-service/config/ /app/config/

# Set the working directory
WORKDIR /app

# Run the application
CMD ["./lgm8-notification-service"]