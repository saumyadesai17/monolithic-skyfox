# Docker Database Setup Guide

## Problem Summary

When running the seed script `sh scripts/seedShowData.sh`, you encountered connection errors because:

1. **Local PostgreSQL Conflict**: You have local PostgreSQL instances running on your Mac (Postgres.app and Homebrew PostgreSQL@14) listening on port 5432
2. **Role Mismatch**: The local PostgreSQL doesn't have the "bookingengine" role that exists in your Docker container
3. **Connection Target**: The `psql` command from your host machine connects to the local PostgreSQL instead of the Docker container

## Solution

Two new scripts have been created to work with your Docker PostgreSQL container:

### 1. Run Migrations: `scripts/run-migrations-docker.sh`

This script runs all database migrations inside the Docker container.

**Usage:**
```bash
cd go-skyfox-backend-base
sh scripts/run-migrations-docker.sh
```

**What it does:**
- Copies migration SQL files to the Docker container
- Executes all migrations in order (slot, show, customer, booking, user tables)
- Runs entirely inside the Docker container to avoid local PostgreSQL conflicts

### 2. Seed Data: `scripts/seedShowData-docker.sh`

This script seeds the database with test data inside the Docker container.

**Usage:**
```bash
cd go-skyfox-backend-base
sh scripts/seedShowData-docker.sh [START_DATE] [NUM_DAYS]

# Examples:
sh scripts/seedShowData-docker.sh 2026-02-24 3
sh scripts/seedShowData-docker.sh 2026-03-01 7
```

**Parameters:**
- `START_DATE` (optional): Starting date for show data (format: YYYY-MM-DD). Defaults to current date.
- `NUM_DAYS` (optional): Number of days to seed. Defaults to 21.

**What it does:**
- Truncates existing data in booking, show, and slot tables
- Seeds 4 time slots (morning, afternoon, evening, late night)
- Seeds random movie shows for the specified number of days
- Runs entirely inside the Docker container

## Environment Variables

The scripts use these environment variables with defaults:

```bash
POSTGRES_CONTAINER=postgresdb    # Docker container name
DB_NAME=bookingengine            # Database name
POSTGRES_USERNAME=bookingengine  # Database user
POSTGRES_PASSWORD=postgres       # Database password
```

You can override them if needed:
```bash
POSTGRES_CONTAINER=mydb sh scripts/seedShowData-docker.sh
```

## Complete Setup Workflow

1. **Ensure Docker PostgreSQL is running:**
   ```bash
   docker ps | grep postgres
   ```

2. **Run migrations (first time only):**
   ```bash
   sh scripts/run-migrations-docker.sh
   ```

3. **Seed the database:**
   ```bash
   sh scripts/seedShowData-docker.sh 2026-02-24 3
   ```

4. **Verify the data:**
   ```bash
   docker exec postgresdb psql -U bookingengine -d bookingengine -c "SELECT COUNT(*) FROM slot; SELECT COUNT(*) FROM show;"
   ```

## Why Not Use the Original Script?

The original `scripts/seedShowData.sh` script is designed to run from the host machine and connect to PostgreSQL via `localhost:5432`. However:

- Your Mac has local PostgreSQL instances running on port 5432
- The `psql` command connects to the local instance instead of Docker
- The local instance doesn't have the "bookingengine" role

**Options to use the original script:**
1. Stop local PostgreSQL services (not recommended if you need them)
2. Change Docker PostgreSQL to use a different port
3. Use the new Docker-based scripts (recommended)

## Troubleshooting

### Check if local PostgreSQL is running:
```bash
ps aux | grep postgres | grep -v grep
```

### Check Docker container status:
```bash
docker ps
docker logs postgresdb
```

### Connect to database manually:
```bash
docker exec -it postgresdb psql -U bookingengine -d bookingengine
```

### View tables:
```sql
\dt
SELECT * FROM slot;
SELECT COUNT(*) FROM show;
```

## Results

After running the scripts successfully:
- **4 slots** created (morning, afternoon, evening, late night)
- **84 shows** created (4 slots × 3 days × 7 shows per day)
- All data properly seeded in the Docker PostgreSQL container
