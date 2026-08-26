# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o notification-service ./cmd/main.go

# Run stage
FROM alpine:3.19
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/notification-service .

EXPOSE 8085
CMD ["./notification-service"]
