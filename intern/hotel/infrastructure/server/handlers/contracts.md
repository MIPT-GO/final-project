## Hotel Service

### Get All Hotels

```
GET /v1/hotel/all/
```

Status Codes:

```
200 OK
```

```
500 Internal Service Error 
```

Response Body:

```json
{
   "hotels": [str, ...]
}
```

### Get All Rooms In Hotel

```
GET /v1/hotel/{hotel: str}/rooms/
```

Status Codes:

```
200 OK
```

```
404 Hotel Not Found (если отель не найден)
```

```
500 Internal Service Error 
```

Response Body:

```json
{
   "numbers": [int, ...]
}
```

### Create Hotel

Request:

```
POST /v1/hotel/
```

Request Body:

```json
{ 
    "name": string,
    "email": string
}
```

Status Codes:

```
201 Created
```

```
400 Bad Request (неверный формат запроса)
```

```
409 Hotel Already Exists
```

```
500 Internal Service Error
```

No response body

### Get One Hotel

Request:

```
GET /v1/hotel/{hotel: str}/
```

Status Codes:

```
200 OK
```

```
404 Hotel Not Found
```

```
500 Internal Service Error
```

Response Body:

```json
{ 
    "name": string,
    "email": string
}
```

### Update Room

Request:

```
PUT /v1/hotel/{hotel: str}/room/
```

Request body:

```json
{
   "number": int,
   "new_cost": float32,
   "user_email": string
}
```

Status Codes:

```
200 OK
```

```
400 Bad Request (неверный формат запроса)
```

```
403 Forbidden (Permission Denied)
```

```
404 Hotel Not Found (или Room Not Found)
```

```
500 Internal Service Error
```

No response body