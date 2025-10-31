# Stage 1: build
FROM golang:1.24.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o subchecker ./cmd/app

# Stage 2: runtime
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/subchecker .
COPY .env .
COPY database/migrations ./database/migrations
EXPOSE 8080
CMD ["./subchecker"]