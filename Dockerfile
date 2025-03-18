# Use the official Go image as the base image
FROM golang:latest

# Set the working directory for the service
WORKDIR /art-service

# Copy `go.mod` and `go.sum` files from the monorepo root directory
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY pkg ./pkg
# COPY .env .


# Change to the service directory and build the application
RUN go build -o art-service ./cmd

# Command to run the application
CMD ["./art-service"]
