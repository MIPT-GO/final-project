# Notification Service

## Environment

- `HOST` — host for metrics bind (metrics: `http://{HOST}:8084/metrics`)
- `LOG_LEVEL` — `debug|info|warning|error`
- `KAFKA_BROKERS` — brokers list (example: `localhost:9092`)
- `KAFKA_TOPIC` — default: `notifications`
- `KAFKA_GROUP_ID` — default: `notification-svc`
- `SMTP_HOST` — default: `smtp.gmail.com`
- `SMTP_PORT` — default: `587`
- `SMTP_USERNAME`
- `SMTP_PASSWORD` — Gmail app password
- `SMTP_FROM` — default: `SMTP_USERNAME`

## Build & Run

Using `go run`:

```bash
go run ./cmd/notification --env path/to/file.env --level info
```

## Kafka message

Topic: `notifications`

Value (JSON):

```json
{
  "event": "booking_created",
  "email": "client@example.com",
  "owner_email": "owner@example.com",
  "hotel": "Hilton",
  "number": 101,
  "start": "2025-12-15T12:00:00Z",
  "end": "2025-12-20T12:00:00Z"
}
```
