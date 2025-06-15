# Watcher

Watches a directory for new folders and publishes events to Kafka.

```
GO111MODULE=on go run ./cmd/watcher
```

Environment variables:
- `WATCH_DIR`  (default `./watch`)
- `KAFKA_TOPIC` (default `folders`)
- `KAFKA_BROKER` (default `localhost:9092`)
