# Build stage
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o watcher ./cmd/watcher

# Final stage
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/watcher ./watcher
CMD ["/app/watcher"]
