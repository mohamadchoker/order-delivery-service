#!/usr/bin/env bash
set -e

echo "🐳 Starting development entrypoint..."

# Construct DB URL dynamically
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

# Wait for Postgres
echo "🗄️  Waiting for Postgres at $DB_HOST:$DB_PORT..."
until PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c '\q' 2>/dev/null; do
  echo "⏳ Postgres is unavailable - sleeping"
  sleep 2
done

echo "✅ Database is ready!"

# Run migrations (skip if migrate not available or if no migrations exist)
if [ -d "./migrations" ] && [ "$(ls -A ./migrations 2>/dev/null)" ]; then
    if command -v migrate >/dev/null 2>&1; then
        echo "🚀 Running database migrations..."
        # Use migrate with explicit postgres:// scheme
        migrate -path ./migrations -database "$DB_URL" up 2>&1 || {
            echo "⚠️  Migration failed or already up to date. Continuing..."
        }
    else
        echo "⚠️  'migrate' command not found. Install with: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
        echo "⚠️  Skipping migrations for now..."
    fi
else
    echo "ℹ️  No migrations found in ./migrations directory"
fi

# Start hot-reloading with Air
# Use AIR_CONFIG env var if set, otherwise default to .air.toml
AIR_CONFIG_FILE="${AIR_CONFIG:-.air.toml}"
echo "🚀 Starting Air for hot-reloading with config: $AIR_CONFIG_FILE"
exec air -c "$AIR_CONFIG_FILE"
