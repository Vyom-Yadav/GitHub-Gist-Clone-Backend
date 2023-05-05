#!/usr/bin/env bash

echo "enabling pg_trgm on database $POSTGRES_DB"
psql -U "$POSTGRES_USER" --dbname="$POSTGRES_DB" <<-'EOSQL'
  CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
EOSQL
echo "finished with exit code $?"
