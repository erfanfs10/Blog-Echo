# Use a minimal base image for Go apps
FROM docker.io/library/golang:1.23.0-alpine3.20

# Set the working directory inside the container
WORKDIR /app

# Copy the source code
COPY . .

# Download dependencies
RUN go mod download

# Build the app
RUN go build -o server

# Expose the port your app uses
EXPOSE 8000

# Command to run the binary
CMD ["./server"]
