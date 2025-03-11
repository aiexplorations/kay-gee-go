FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/frontend ./cmd/frontend

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/frontend /app/frontend
COPY public /app/public

EXPOSE 8080

ENTRYPOINT ["/app/frontend"] 