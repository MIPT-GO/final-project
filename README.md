# Как запустить проект
Нужно создать .env файлы из примеров в config и положить рядом.

К примеру:
./config/db/booking_db.env.template -> ./config/db/booking_db.env
и т.д.

После того как все .env файлы созданы, проект нужно запустить через 
```bash
docker compose -f ./deploy/docker-compose.yml up -d
```

# Troubleshooting
Если ваша система не linux, то в файле ./deploy/compose/observability.yml
нужно убрать все volumes, кроме
- ${PWD}/config/config.alloy:/etc/alloy/config.alloy:ro

# Endpoints

## Payment System

### Create Payment Link

```
POST /{router-prefix}/link
```

Body:
```json
{
    "amount": "string",
    "webhook": "string",
    "message": "string"
}
```

Response may be:

```
200 OK
```

```
422 Invalid request body
```

```
500 Internal Service Error
```

If code is 200. System return json:

```json
{
    "link": "string-url"
}
```

### Receiving Payment Result

When the user clicks on the link and clicks pay or time`s up, the system will send the following request to the address specified in the webhook field:

```json
{
    "status": "string"
}
```

Status meaning:

- `ok` - payment is successful
- `timeout` - the waiting time has expired, timeout

## Hotel Service

### Environment 

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

### Endpoints

- `GET /v1/hotel/all/` — list all hotels.
- `POST /v1/hotel/` — create a hotel (body contains hotel fields).
- `GET /v1/hotel/{hotel}/` — get hotel details (name and email).
- `GET /v1/hotel/{hotel}/rooms/` — list all room numbers in the hotel.
- `GET /v1/hotel/{hotel}/room/{number}` — get room cost.
- `PUT /v1/hotel/{hotel}/room/` — update room cost (body contains room number and new cost).
- `POST /v1/hotel/{hotel}/room/` — create a room (body contains room fields).

Подробнее: [Hotel Service Readme](../hotel/README.md)

