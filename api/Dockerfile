# Build stage
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o api ./cmd/api

# Final stage
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/api ./api
COPY --from=builder /app/public ./public
EXPOSE 8080
ENTRYPOINT ["/app/api"]
