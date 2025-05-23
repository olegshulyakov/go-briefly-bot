# Stage 1: Build with Maven
FROM maven:3.9-eclipse-temurin-21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy Maven pom.xml
COPY pom.xml .

# Download dependencies
RUN mvn dependency:go-offline

# Copy the rest of the application code
COPY src ./src

# Build the application
RUN mvn package -DskipTests

# Stage 2: Create a lightweight runtime image
FROM eclipse-temurin:21-jre-alpine

# Install yt-dlp and Python (required for yt-dlp)
RUN apk add --no-cache yt-dlp

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/target/*.jar youtube-reteller.jar

# Ports and volumes
EXPOSE 8080

# Set the entrypoint to run the application
ENTRYPOINT ["java", "-jar", "youtube-reteller.jar"]