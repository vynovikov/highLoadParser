#!/bin/sh
echo "starting..."

/dlv --accept-multiclient --continue --headless=true --listen=:40000 --api-version=2 --check-go-version=false exec /app/main