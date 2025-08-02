# Use the official Golang image
FROM golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Use a smaller image for the final container
FROM alpine:latest

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/app .

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./app"]
