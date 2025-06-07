# Go Streaming Monorepo

This repository contains three Go projects:

1. **api** - File system API using [Fiber](https://gofiber.io/) to stream videos.
2. **watcher** - A program that watches a Linux directory and publishes an event to Kafka when a new folder is detected.
3. **processor** - A Kafka consumer that processes videos using `ffmpeg` and writes the results to an output directory that the API can stream.

Each project is maintained as an independent Go module.
