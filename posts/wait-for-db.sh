#!/bin/sh

set -e

host="$DB_HOST"
port="$DB_PORT"

# Esperar hasta que la base de datos esté disponible
until nc -z "$host" "$port"; do
  echo "Waiting for database at $host:$port..."
  sleep 1
done

echo "Database is up!"
exec "$@"