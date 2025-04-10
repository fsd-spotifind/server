#!/bin/sh

# Wait for the database to be ready
echo "Waiting for database to be ready..."
while ! nc -z psql_bp 5432; do
  sleep 1
done

# Run migrations
echo "Running database migrations..."
go run github.com/steebchen/prisma-client-go migrate deploy

# Start the application
echo "Starting application..."
./main 