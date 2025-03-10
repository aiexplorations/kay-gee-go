FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o enricher ./cmd/kg-enricher

# Use a minimal alpine image for the final image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/enricher /app/enricher

# Create cache directory
RUN mkdir -p /app/cache/llm

# Make sure the binary is executable
RUN chmod +x /app/enricher

# Set the entrypoint
ENTRYPOINT ["/app/enricher"] 