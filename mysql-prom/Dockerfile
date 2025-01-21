# Stage 1: Build the Go binary
FROM golang:1.23 as builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o app main.go

# # Stage 2: Create the final container
# FROM alpine:latest

# # Install CA certificates for HTTPS
# RUN apk --no-cache add ca-certificates

# # Set the working directory
# WORKDIR /root/

# # Copy the compiled binary from the builder stage
# COPY --from=builder /app/app .

# Expose the metrics port
EXPOSE 8080

# Command to run the application
CMD ["/app/app"]
