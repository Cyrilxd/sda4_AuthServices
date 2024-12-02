# Use the official Go image as the base image
FROM golang:1.23.3 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth-service ./cmd/main.go

# Use a minimal base image for the final container
FROM alpine:latest

# Set the working directory for the final image
WORKDIR /root/

# Copy the compiled binary from the builder
COPY --from=builder /app/auth-service .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the binary
CMD ["./auth-service"]
