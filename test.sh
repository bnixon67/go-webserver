#!/usr/bin/env bash

BASE_URL="http://localhost:8080"
OPTS="-s -i"

endpoints=(
    "/"
    "/hello"
    "/headers"
    "/request"
    "/remote"
)

for endpoint in "${endpoints[@]}"; do
    echo curl $OPTS "$BASE_URL$endpoint"
    curl $OPTS "$BASE_URL$endpoint"
    echo "---------------------------------------------------------------------"
done
