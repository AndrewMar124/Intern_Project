# Use the official Golang image to build the app
FROM golang as builder

WORKDIR /app
# Copy the source from the current directory to the Working Directory inside the container
# param 1: local package go
# param 2: where we want to copy it
COPY . .

CMD go run main.go