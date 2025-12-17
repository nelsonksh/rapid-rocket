# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod ./
# COPY go.sum ./ # Uncomment if you have dependencies later

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=0 creates a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Run Stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests (Andamio API)
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy static assets and templates
COPY --from=builder /app/views ./views
COPY --from=builder /app/assets ./assets

# Expose the port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
