#!/bin/bash

# This script is used to run database migrations.

# Exit immediately if a command exits with a non-zero status.
set -e

# Define the database connection parameters
DB_USER="your_db_user"
DB_PASSWORD="your_db_password"
DB_NAME="your_db_name"
DB_HOST="localhost"
DB_PORT="5432"

# Run the migration SQL file
psql -U $DB_USER -d $DB_NAME -h $DB_HOST -p $DB_PORT -f ./internal/database/migrations/001_create_users_table.sql

echo "Database migration completed successfully."