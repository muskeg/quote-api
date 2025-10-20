#!/bin/sh
set -e

# Create data directory if it doesn't exist
mkdir -p /app/data

# Copy default config if missing
if [ ! -f /app/data/config.yaml ]; then
  echo "Initializing default config.yaml..."
  cp /app/config-default.yaml /app/data/config.yaml
fi

# Copy default quotes if missing
if [ ! -f /app/data/quotes.json ]; then
  echo "Initializing default quotes.json..."
  cp /app/quotes-default.json /app/data/quotes.json
fi

# Run the main application
exec /app/quote-api