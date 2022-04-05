#!/bin/bash
psql --username postgres --dbname file_storage <<-EOSQL
CREATE TABLE IF NOT EXISTS files (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    file BYTEA
);
CREATE TABLE IF NOT EXISTS records (
    file_id VARCHAR(255) REFERENCES files(id),
    username VARCHAR(255),
    created_at VARCHAR(255),
    processed BOOLEAN
);
EOSQL