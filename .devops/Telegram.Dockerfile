# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder

ENV GOROOT /usr/local/go

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o go-briefly-telegram ./cli/telegram

# Stage 2: Create a lightweight runtime image
FROM alpine:latest

# Install yt-dlp and Python (required for yt-dlp)
RUN apk add --no-cache yt-dlp ffmpeg

# Set the working directory
WORKDIR /app

# Create the data directory
RUN mkdir -p /app/data

# Add volume
VOLUME ["/app/data"]

# Copy the built binary from the builder stage
COPY --from=builder /app/go-briefly-telegram /usr/bin/

# Copy the locales directory
COPY --from=builder /app/locales /app/locales

# Set the entrypoint to run the application
ENTRYPOINT ["go-briefly-telegram"]