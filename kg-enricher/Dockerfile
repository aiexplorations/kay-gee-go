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

# Build the enricher
RUN go build -o /kg-enricher ./cmd/kg-enricher

# Expose the API port
EXPOSE 5001

# No entrypoint, will be set in docker-compose.yml 