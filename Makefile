# Database configuration
PG_USER = postgres
PG_PASSWORD = 12345678
PG_DBNAME = beego-presence-api
PG_HOST = localhost
PG_PORT = 5432

# Export password for non-interactive use
export PGPASSWORD=$(PG_PASSWORD)

.PHONY: reset-db drop-db create-db

# Default target
reset-db: drop-db create-db
	@echo "Database has been refreshed!"

# Drop the database
drop-db:
	@echo "Dropping database $(PG_DBNAME)..."
	@psql -U $(PG_USER) -h $(PG_HOST) -p $(PG_PORT) -c "DROP DATABASE IF EXISTS \"$(PG_DBNAME)\";" || true

# Create the database
create-db:
	@echo "Creating database $(PG_DBNAME)..."
	@psql -U $(PG_USER) -h $(PG_HOST) -p $(PG_PORT) -c "CREATE DATABASE \"$(PG_DBNAME)\";" || true
