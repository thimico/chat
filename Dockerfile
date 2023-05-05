FROM golang:1.17

WORKDIR /app

COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code into the container
COPY . .

# Build the application
RUN go build -o main .

# Expose the port the app runs on
EXPOSE 8080

# Run the binary
CMD ["/app/main"]
