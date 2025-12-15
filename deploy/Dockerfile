# syntax=docker/dockerfile:1
FROM golang:1.25-alpine AS builder

LABEL maintainer="Zherdev Egor"
LABEL vendor="Go final project team"

ARG APP="cmd"

WORKDIR /src

RUN apk update && apk add --no-cache \
    ca-certificates  \
    git  \
    make

COPY go.mod go.sum ./
RUN go mod download && go mod verify

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -a \
    -installsuffix cgo \
    -trimpath \
    -o /src/bin/app \
    "./cmd/${APP}"

FROM alpine:3.21 AS runtime

ARG UID=1000
ARG GID=1000

RUN apk update && apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    postgresql-client \
    bash  \
    && curl -fsSL https://raw.githubusercontent.com/pressly/goose/v3.20.0/install.sh | bash

WORKDIR /app

RUN addgroup -g "${GID}" -S gopher && \
    adduser -S gopher -u "${UID}" -G gopher && \
    mkdir -p /app && \
    chown -R gopher:gopher /app

COPY --from=builder --chown=gopher:gopher /src/bin/app .

COPY --chown=gopher:gopher deploy/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

USER gopher:gopher

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["./app"]
