# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS builder

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

WORKDIR /app

SHELL ["/bin/bash", "-eo", "pipefail", "-c"]

RUN apk update && apk add --no-cache \
    ca-certificates \
    tzdata \
    curl

RUN addgroup -g "${GID}" -S gohper && \
    adduser -d '/app' -S gopher -u "${UID}" -G gopher && \
    chown gopher:gopher -R /app

COPY --from=builder --chown=gopher:gopher /src/bin/app .

USER gopher

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

EXPOSE 8080

ENTRYPOINT ["./app"]
CMD ["serve"]
