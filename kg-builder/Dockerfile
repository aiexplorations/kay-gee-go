FROM golang:1.20-alpine

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the API server
RUN go build -o /app/api-server api-server.go

# Build the builder
RUN go build -o /kg-builder ./cmd/kg-builder

# Expose the API port
EXPOSE 5000

# No entrypoint, will be set in docker-compose.yml
