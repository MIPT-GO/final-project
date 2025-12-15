# Hotel Service

## Description

Сервис для управления сущностями отелей и комнат, предоставляющий RESTful API и использующий структурированное
логирование.

## Launch

Исполняемый файл сервиса требует корректно установленных переменных окружения для конфигурации.

### Environment Variables

Система использует следующие переменные окружения, в основном для подключения к базе данных и настройки сервера:

| Variable            | Default Value |
| :------------------ | :------------ |
| `HOST`              | `0.0.0.0`     |
| `PORT`              | `8080`        |
| `DB_HOST`           | `localhost`   |
| `DB_PORT`           | `5432`        |
| `DB_USER`           | `postgres`    |
| `DB_PASSWORD`       | `secret`      |
| `DB_DATABASE`       | `hotel_db`    |
| `READ_TIMEOUT_SEC`  | `5`           |
| `WRITE_TIMEOUT_SEC` | `10`          |
| `IDLЕ_TIMEOUT_SEC`  | `15`          |

## API Usage

Все конечные точки имеют префикс `/v1/hotel`.

### Get All Hotels

Получает список всех зарегистрированных названий отелей.

```
GET /v1/hotel/all/
```

Status Codes:

```
200 OK
500 Internal Service Error 
```

Response Body:

```json
{
  "hotels": [
    "string",
    ...
  ]
}
```

### Create Hotel

Создает новую запись об отеле.

```
POST /v1/hotel/
```

Request Body:

```json
{
  "name": "string",
  "email": "string"
}
```

Status Codes:

```
201 Created
400 Bad Request (Неверный формат)
409 Conflict (Hotel Already Exists)
500 Internal Service Error
```

No response body.

### Get One Hotel

Получает детали (имя и email) конкретного отеля.

```
GET /v1/hotel/{hotel: str}/
```

Status Codes:

```
200 OK
404 Not Found (Hotel Not Found)
500 Internal Service Error
```

Response Body:

```json
{
  "name": "string",
  "email": "string"
}
```

### Get All Rooms In Hotel

Получает список номеров всех комнат в указанном отеле.

```
GET /v1/hotel/{hotel: str}/rooms/
```

Status Codes:

```
200 OK
500 Internal Service Error
```

Response Body:

```json
{
  "numbers": [
    int,
    ...
  ]
}
```

### Get Room Cost

Получает текущую стоимость конкретной комнаты в отеле.

```
GET /v1/hotel/{hotel: str}/room/{number: int}
```

Status Codes:

```
200 OK
400 Bad Request (Invalid room number format)
404 Not Found (Room Not Found)
500 Internal Service Error
```

Response Body:

```json
{
  "cost": 150.50
}
```

### Update Room Cost

Обновляет цену комнаты. Требуется `user_email` для проверки разрешений (должен быть владельцем отеля).

```
PUT /v1/hotel/{hotel: str}/room/
```

Request body:

```json
{
  "number": 101,
  "new_cost": 150.50,
  "user_email": "owner@hotel.com"
}
```

Status Codes:

```
200 OK
400 Bad Request (Invalid format)
403 Forbidden (Permission Denied)
404 Not Found (Hotel or Room Not Found)
500 Internal Service Error
```

No response body.


```
POST /v1/hotel/{hotel: str}/room/
```

Request body:

```json
{
  "number": 101,
  "cost": 150.50,
  "user_email": "owner@hotel.com"
}
```

Status Codes:

```
201 Created
400 Bad Request (Неверный формат)
409 Conflict (Room Already Exists)
500 Internal Service Error
```

No response body.

