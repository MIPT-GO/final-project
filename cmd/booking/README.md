# Booking Service

## Environment (env/booking.dev.env)
- `DB_HOST` — Postgres host (example: `localhost`)
- `DB_PORT` — Postgres port (example: `5432`)
- `DB_USER` — Postgres user (example: `postgres`)
- `DB_PASSWORD` — Postgres password (example: `postgres`)
- `DB_DATABASE` — Postgres database name (example: `booking`)

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


# Hotel Service

## Environment 

- `DB_HOST` — Postgres host (example: `localhost`)
- `DB_PORT` — Postgres port (example: `5432`)
- `DB_USER` — Postgres user (example: `postgres`)
- `DB_PASSWORD` — Postgres password (example: `secret`)
- `DB_DATABASE` — Postgres database name (example: `hotel_db`)

- `SERVER_HOST` — Service host (example: `0.0.0.0`)
- `SERVER_PORT` — Service port (example: `8080`)

- `READ_TIMEOUT_SEC` — HTTP read timeout in seconds (example: `5`)
- `WRITE_TIMEOUT_SEC` — HTTP write timeout in seconds (example: `10`)
- `IDLE_TIMEOUT_SEC` — HTTP idle timeout in seconds (example: `15`)

## Endpoints

- `GET /v1/hotel/all/` — list all hotels.
- `POST /v1/hotel/` — create a hotel (body contains hotel fields).
- `GET /v1/hotel/{hotel}/` — get hotel details (name and email).
- `GET /v1/hotel/{hotel}/rooms/` — list all room numbers in the hotel.
- `GET /v1/hotel/{hotel}/room/{number}` — get room cost.
- `PUT /v1/hotel/{hotel}/room/` — update room cost (body contains room number and new cost).
- `POST /v1/hotel/{hotel}/room/` — create a room (body contains room fields).

Подробнее: [Hotel Service Readme](../hotel/README.md)

## Migrations

* Migrations are located in the `migrations/` directory.


## Build & Run

Using `go run`:
```bash
go run ./cmd/booking -loglevel `debug|info|warn|error` -env 
```
