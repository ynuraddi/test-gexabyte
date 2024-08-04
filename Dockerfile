FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o gexabyte ./cmd/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/gexabyte .
COPY --from=builder /app/internal/repository/postgres/migrations /app/migrations

LABEL maintainers = "ynuraddi"
LABEL version = "1.0"

EXPOSE 8080

CMD ["./gexabyte"]