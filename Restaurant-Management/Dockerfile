# Use official Golang image as the base
FROM golang:1.23

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port Gin uses (default 8080)
EXPOSE 8080

# Command to run the app
CMD ["/app/main"]
