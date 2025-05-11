# Build stage
FROM golang:alpine AS builder
WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o server main.go

# Run stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .

EXPOSE 5000
CMD ["./server"]
