echo "Executing load-db-extensions"
psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
\\q
echo "OK"