# Use the official Golang base image
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the local package files to the container's working directory
COPY . .

# Build the Go application
RUN go build -o main ./cmd/

# Use the official smaller linux base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy only the necessary files from the builder stage
COPY --from=builder /app/main ./

# Expose the port the application runs on
EXPOSE 8080

# Command to run the executable
RUN ls
CMD ["./main", "-env=dev"]