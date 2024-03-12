FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN go build -v -o bashapi ./cmd/...

FROM alpine:latest

COPY --from=builder /app/bashapi .

CMD ["./bashapi"]
