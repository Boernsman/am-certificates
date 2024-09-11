# Use the official Golang image to build the Go backend
FROM golang:1.20-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main .

# Use a lightweight image for running the app
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built Go application from the builder stage
COPY --from=builder /app/main .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the Go app
CMD ["./main"]
