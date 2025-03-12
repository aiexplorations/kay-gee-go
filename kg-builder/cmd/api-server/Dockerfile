FROM golang:1.20-alpine

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the API server
RUN go build -o api-server ./cmd/api-server

# Expose the API port
EXPOSE 5000

# Run the API server
CMD ["/app/api-server"] 