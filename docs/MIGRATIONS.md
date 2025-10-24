# Database Migrations Guide

## Overview

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database schema management.

## Migration Strategies

### ✅ Recommended: Docker-Based Migrations

Use Docker for migrations to ensure consistency across all environments:

```bash
# Check current migration version
make migrate-status-docker

# Apply all pending migrations
make migrate-up-docker

# Rollback last migration
make migrate-down-docker
```

**Why Docker?**
- ✅ Works with `make dev-up` or `make dev-debug`
- ✅ Same environment as your application
- ✅ Network access to Docker PostgreSQL
- ✅ No local setup required
- ✅ Consistent across all developers
- ✅ Works in CI/CD

**Automatic migrations:** When you start the dev environment with `make dev-up`, migrations run automatically via the entrypoint script.

### Local Migrations (Alternative)

If you have PostgreSQL running locally (outside Docker):

```bash
# Apply migrations
make migrate-up

# Rollback migrations
make migrate-down
```

**Requirements:**
- PostgreSQL running on `localhost:5432`
- Database `order_delivery_db` exists
- `~/go/bin` in your PATH
- `migrate` tool installed with PostgreSQL support

## Creating New Migrations

```bash
# Create a new migration
make migrate-create NAME=add_user_table

# This creates two files:
# migrations/000002_add_user_table.up.sql   - Forward migration
# migrations/000002_add_user_table.down.sql - Rollback migration
```

**Migration naming conventions:**
- Use snake_case
- Be descriptive: `add_status_column`, not `update_table`
- Format: `000001_description.up.sql` and `000001_description.down.sql`

## Writing Migrations

### Up Migration (Forward)

File: `migrations/000002_add_user_table.up.sql`

```sql
-- Add user table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Add index
CREATE INDEX idx_users_email ON users(email);

-- Add comments
COMMENT ON TABLE users IS 'User accounts';
COMMENT ON COLUMN users.email IS 'User email address (unique)';
```

### Down Migration (Rollback)

File: `migrations/000002_add_user_table.down.sql`

```sql
-- Rollback user table
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
```

### Best Practices

1. **Always use IF EXISTS/IF NOT EXISTS**
   ```sql
   CREATE TABLE IF NOT EXISTS ...
   DROP TABLE IF EXISTS ...
   ```

2. **Make migrations reversible**
   - Every `up` migration should have a corresponding `down`
   - Test both directions

3. **One logical change per migration**
   - Don't combine unrelated changes
   - Easier to debug and rollback

4. **Use transactions implicitly**
   - PostgreSQL wraps DDL in transactions automatically
   - But be careful with large data migrations

5. **Add comments**
   ```sql
   COMMENT ON TABLE delivery_assignments IS 'Tracks delivery lifecycle';
   ```

6. **Handle existing data**
   ```sql
   -- Safe: Add column with default
   ALTER TABLE users ADD COLUMN status VARCHAR(20) DEFAULT 'active';

   -- Risky: Add NOT NULL without default (will fail if table has data)
   -- ALTER TABLE users ADD COLUMN status VARCHAR(20) NOT NULL;
   ```

## Migration Workflow

### Adding a New Feature

1. **Create migration**
   ```bash
   make migrate-create NAME=add_notifications_table
   ```

2. **Write SQL**
   - Edit `000N_add_notifications_table.up.sql`
   - Edit `000N_add_notifications_table.down.sql`

3. **Test locally**
   ```bash
   # Apply migration
   make migrate-up-docker

   # Verify it worked
   make migrate-status-docker

   # Test rollback
   make migrate-down-docker

   # Re-apply
   make migrate-up-docker
   ```

4. **Commit migration files**
   ```bash
   git add migrations/
   git commit -m "Add notifications table migration"
   ```

5. **Migrations run automatically** when other developers pull and run `make dev-up`

### Fixing a Migration

**If not yet committed:**
```bash
# Rollback
make migrate-down-docker

# Edit the SQL files
vim migrations/000N_your_migration.up.sql

# Re-apply
make migrate-up-docker
```

**If already committed/pushed:**
- **Never edit existing migrations** that others might have run
- Create a new migration to fix the issue
- Example: `make migrate-create NAME=fix_users_table_index`

## Checking Migration Status

```bash
# Via Docker (recommended)
make migrate-status-docker

# Output example:
# 1  # Current version (migration 000001 applied)
```

## Troubleshooting

### "No running container found"

Start the dev environment first:
```bash
make dev-up
# or
make dev-debug
```

### "Dirty database version"

This happens when a migration fails partway through:

```bash
# Via Docker
docker exec order-delivery-service-dev sh -c 'migrate -path /app/migrations -database "postgres://postgres:postgres@postgres:5432/order_delivery_db?sslmode=disable" force <version>'

# Replace <version> with the version number you want to force
```

### "Migration already exists"

You ran the migration twice. Check current version:
```bash
make migrate-status-docker
```

### Check migration history in database

```bash
# Connect to database
docker exec -it order-delivery-postgres-dev psql -U postgres -d order_delivery_db

# Check schema_migrations table
SELECT * FROM schema_migrations;

# Exit
\q
```

## CI/CD Integration

Migrations run automatically in Docker containers, so your CI/CD pipeline should:

1. Build Docker image
2. Start containers with `docker-compose up`
3. Migrations run via entrypoint
4. Run tests
5. Deploy

No special migration step needed!

## Advanced Usage

### Run specific migration version

```bash
# Via Docker
docker exec order-delivery-service-dev sh -c 'migrate -path /app/migrations -database "postgres://postgres:postgres@postgres:5432/order_delivery_db?sslmode=disable" goto <version>'
```

### Check available commands

```bash
docker exec order-delivery-service-dev migrate --help
```

## Common Patterns

### Adding a Column

```sql
-- Up
ALTER TABLE delivery_assignments
ADD COLUMN priority VARCHAR(20) DEFAULT 'normal';

-- Down
ALTER TABLE delivery_assignments
DROP COLUMN IF EXISTS priority;
```

### Adding an Index

```sql
-- Up
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_delivery_status
ON delivery_assignments(status)
WHERE deleted_at IS NULL;

-- Down
DROP INDEX CONCURRENTLY IF EXISTS idx_delivery_status;
```

### Modifying a Column

```sql
-- Up: Change column type
ALTER TABLE delivery_assignments
ALTER COLUMN driver_id TYPE VARCHAR(100);

-- Down: Restore original type
ALTER TABLE delivery_assignments
ALTER COLUMN driver_id TYPE VARCHAR(50);
```

### Adding Foreign Key

```sql
-- Up
ALTER TABLE delivery_assignments
ADD CONSTRAINT fk_delivery_driver
FOREIGN KEY (driver_id)
REFERENCES drivers(id)
ON DELETE SET NULL;

-- Down
ALTER TABLE delivery_assignments
DROP CONSTRAINT IF EXISTS fk_delivery_driver;
```

## Migration File Location

All migrations live in:
```
migrations/
├── 000001_initial_schema.up.sql
├── 000001_initial_schema.down.sql
├── 000002_add_status_index.up.sql
└── 000002_add_status_index.down.sql
```

Never delete migration files that have been applied in production!
