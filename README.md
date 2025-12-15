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
