# Booking Service

## Environment (env/booking.dev.env)
- `PG_HOST` — Postgres host (example: `localhost`)
- `PG_PORT` — Postgres port (example: `5432`)
- `PG_USER` — Postgres user (example: `postgres`)
- `PG_PASSWORD` — Postgres password (example: `postgres`)
- `PG_DATABASE` — Postgres database name (example: `booking`)

- `HOTEL_SERVICE_HOST` — Hotel service base URL (include scheme, e.g. `http://localhost`)
- `HOTEL_SERVICE_PORT` — Hotel service port (example: `8081`)
- `HOTEL_SERVICE_TIMEOUT` — HTTP timeout (example: `5s`)

- `PAYMENT_SERVICE_HOST` — Payment service base URL (include scheme, e.g. `http://localhost`)
- `PAYMENT_SERVICE_PORT` — Payment service port (example: `8081`)
- `PAYMENT_SERVICE_TIMEOUT` — HTTP timeout for payment service (example: `5s`)

- `SERVER_HOST` — Public host for this service (must include scheme, e.g. `http://localhost`)
- `SERVER_PORT` — Port passed to `http.Server` (format `:8080`)

Example dev file: `env/booking.dev.env`

## Endpoints
- `POST /v1/booking/` — create a reservation (body contains reservation fields). Returns JSON `{ "payment_link": "..." }`.
- `GET /v1/booking/available/{email}/` — list reservations by email.
- `GET /v1/booking/hotel/{hotel}/` — list reservations by hotel.

## Migrations
- Migrations are in the `migrations/` folder.

## Build & Run

Using `go run`:
```bash
go run ./cmd/booking -loglevel `debug|info|warn|error` -env 
```
