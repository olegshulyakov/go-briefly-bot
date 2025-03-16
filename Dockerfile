# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o telegram-youtube-reteller .

# Stage 2: Create a lightweight runtime image
FROM alpine:latest

# Install yt-dlp and Python (required for yt-dlp)
RUN apk add --no-cache yt-dlp

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/telegram-youtube-reteller .

# Set the entrypoint to run the application
ENTRYPOINT ["./telegram-youtube-reteller"]