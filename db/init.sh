#!/bin/bash
psql --username postgres --dbname file_storage <<-EOSQL
CREATE TABLE IF NOT EXISTS files (name VARCHAR(255) PRIMARY KEY, file BYTEA);
EOSQL