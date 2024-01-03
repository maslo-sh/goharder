# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory to the app directory
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Build the Go app
RUN go build -o main ./cmd

# Command to run the executable
CMD ["./main"]
