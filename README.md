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
