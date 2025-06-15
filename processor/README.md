# Processor

Consumes folder events from Kafka and runs `ffmpeg` on video files inside each folder.

```
GO111MODULE=on go run ./cmd/processor
```

Environment variables:
- `KAFKA_TOPIC` (default `folders`)
- `KAFKA_BROKER` (default `localhost:9092`)
- `OUT_DIR` (default `./out`)
