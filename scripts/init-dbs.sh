#!/bin/bash
set -e

echo "í´§ Starting database initialization..."

# Create the products database
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE cloudcart_products;
    GRANT ALL PRIVILEGES ON DATABASE cloudcart_products TO cloudcart;
EOSQL

echo "âœ… Database cloudcart_products created successfully"
