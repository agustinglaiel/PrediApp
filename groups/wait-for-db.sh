#!/bin/sh

echo "⏳ Esperando a que MySQL esté listo en $DB_HOST:3306..."

# Espera hasta que pueda conectarse a MySQL
until nc -z "$DB_HOST" 3306; do
  echo "⏱️  Aún no disponible. Esperando..."
  sleep 2
done

echo "✅ MySQL está disponible. Iniciando el servicio..."
exec "$@"
