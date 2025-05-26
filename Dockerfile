# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o logleaf-server ./cmd/server/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/logleaf-server ./
EXPOSE 8080
EXPOSE 8082
EXPOSE 8082
ENV GIN_MODE=release
CMD ["./logleaf-server"]
