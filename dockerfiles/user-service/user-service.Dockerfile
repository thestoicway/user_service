# Start from the official Golang base image for building the application
FROM golang:1.21.4 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.work file if you are using workspace mode in Go
COPY go.work .

# Copy the entire user_service directory
COPY ./user_service ./user_service

# Copy the custom_errors package
COPY ./custom_errors ./custom_errors

# Download necessary Go modules
RUN cd user_service && go mod download

# Set the working directory to the main application directory
WORKDIR /app/user_service/cmd/user

# Compile the application
RUN go build -o user_service .

# Use a debian slim image for running the application
# Don't use the alpine image as it doesn't have the necessary libraries
FROM debian:stable-slim

# Set the working directory in the new image
WORKDIR /root/

# Copy the compiled application from the builder stage
COPY --from=builder /app/user_service/cmd/user/user_service .

RUN chmod +x user_service

# Command to run the executable
CMD ["./user_service"]