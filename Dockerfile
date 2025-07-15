FROM golang:1.23.6-alpine AS builder
RUN apk add --no-cache gcc musl-dev git
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=1
RUN go build -o bot cmd/bot/main.go

FROM alpine:3.18
RUN apk add --no-cache sqlite ca-certificates
WORKDIR /app
COPY --from=builder /build/bot /app/bot
COPY config /app/config
COPY .env /app/.env
RUN chmod +x /app/bot
EXPOSE 8080
CMD ["./bot"]
