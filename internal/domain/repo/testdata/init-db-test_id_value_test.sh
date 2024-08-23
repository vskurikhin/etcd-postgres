#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE TABLE IF NOT EXISTS test_id_value_test (id BIGSERIAL PRIMARY KEY, value TEXT);
    INSERT INTO test_id_value_test (id, value) VALUES (1, 'test')
EOSQL
