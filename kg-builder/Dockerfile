FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum* ./

RUN go mod download

COPY . .

RUN go build -o /kg-builder ./cmd/kg-builder

CMD ["/kg-builder"]
