FROM golang:1.19-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY src/go.mod src/go.sum ./

# Download dependencies and verify
RUN go mod download && go mod verify

# Copy source code
COPY src/ ./

# Tidy up dependencies
RUN go mod tidy

# Build the application
RUN go build -o kg-frontend .

# Create final image
FROM alpine:3.16

# Install runtime dependencies
RUN apk add --no-cache bash docker-cli

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/kg-frontend /app/kg-frontend

# Copy shell scripts
COPY src/*.sh /app/
RUN chmod +x /app/*.sh

# Copy static files
COPY public/ /app/public/

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080
ENV NEO4J_URI=bolt://neo4j:7687
ENV NEO4J_USER=neo4j
ENV NEO4J_PASSWORD=password

# Run the application
CMD ["/app/kg-frontend"] 