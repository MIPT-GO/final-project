# Payment System

## Description

A simple test payment system based on webhook mechanisms

## Launch

Service executable file has two flags

- `--env` Set env file to load
- `--level` Set logging level

System use this list of environment variables:

- `HOST` 
- `PORT` 
- `REDIS_HOST` 
- `REDIS_PORT`
- `REDIS_DB` 
- `REDIS_USER`
- `REDIS_PASSWORD` 
- `LINK_PREFIX` - The prefix for all returned payment links must match this pattern: `http[s]://{url-to-system}/{router-prefix}/page`
- `ROUTER_PREFIX` - A prefix for all application handlers. By default `/v1/payment`
- `PAYMENT_TIMEOUT` - The maximum number of minutes to pay using the created link
- `STORAGE_TIMEOUT` - The maximum number of minutes that payment information is stored in the system must be strictly longer than `PAYMENT_TIMEOUT`
- `TEMPLATE_DIR` - The path to search for the folder with templates for payment pages. Basic templates located at `intern/payment-system/infrastructure/templates`

You can view an example of the created environment file at `.payment-system.dev.env`

## Using System

### Create Payment Link

To create payment link, send POST request to `http://{url}/{router-prefix}/link`

With this body:

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
