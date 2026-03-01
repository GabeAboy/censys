FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o censys main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

# Copy binary from builder & migration files
COPY --from=builder /app/censys .
COPY --from=builder /app/internal/database/migration ./internal/database/migration

EXPOSE 8080

CMD ["./censys"]
