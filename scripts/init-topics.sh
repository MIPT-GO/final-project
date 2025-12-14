#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

readonly BOOTSTRAP_SERVER="${BOOTSTRAP_SERVER:-kafka:9092}"

create_topic() {
  local name=$1
  local partitions=${2:-3}
  local replication=${3:-1}

  echo "Creating topic: $name"
  kafka-topics --bootstrap-server "$BOOTSTRAP_SERVER" \
    --create --if-not-exists \
    --topic "$name" \
    --partitions "$partitions" \
    --replication-factor "$replication"
}

create_topic "notifications" 1 1

echo "Done"
