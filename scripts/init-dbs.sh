#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE cloudcart_products;
    GRANT ALL PRIVILEGES ON DATABASE cloudcart_products TO $POSTGRES_USER;
EOSQL

echo "✅ Database cloudcart_products created"
