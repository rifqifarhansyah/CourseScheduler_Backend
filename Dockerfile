# Use the official Golang image as the base image
FROM golang:alpine

# Set the working directory in the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project to the container
COPY . .

# Build the Go application
RUN go build -o main .

# Set the command to run the application
CMD ["./main"]
