# Use the official Golang image to build the app
FROM golang as builder
WORKDIR /app
# Copy the source from the current directory to the Working Directory inside the container
COPY . .

CMD go run main.go