# Use the official Golang image to build the app
FROM golang as builder

# working directory instide container
WORKDIR /app

#copy these files
COPY go.mod go.sum ./

# Downlload all dependancies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
# param 1: local package go
# param 2: where we want to copy it
COPY . .

# Build all files in package "main"
Run go build -o main .

# Execute
CMD ["./main"]