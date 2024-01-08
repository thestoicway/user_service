# Start from the official Golang base image for building the application
FROM golang:1.21.4 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download necessary Go modules
RUN go mod download

# Copy the entire user_service directory
COPY . .

# Set the working directory to the main application directory
WORKDIR /app/cmd

# Compile the application
RUN go build -o user_service .

# Use a debian slim image for running the application
# Don't use the alpine image as it doesn't have the necessary libraries
FROM debian:stable-slim

# Set the working directory in the new image
WORKDIR /root/

# Copy the compiled application from the builder stage
COPY --from=builder /app/cmd/user_service .

RUN chmod +x user_service

# Command to run the executable
CMD ["./user_service"]