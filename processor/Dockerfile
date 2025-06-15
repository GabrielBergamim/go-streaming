# Build stage
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o processor ./cmd/processor

# Final stage
FROM alpine:latest
RUN apk add --no-cache ffmpeg
WORKDIR /app
COPY --from=builder /app/processor ./processor
CMD ["/app/processor"]
