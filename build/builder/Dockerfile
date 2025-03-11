FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/builder ./cmd/builder

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/builder /app/builder
COPY config/builder.yaml /app/config.yaml

RUN mkdir -p /app/cache

ENTRYPOINT ["/app/builder"] 