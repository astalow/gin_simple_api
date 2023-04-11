#!/bin/bash

API_ENDPOINT="localhost:8080"
JSON_DATA='{"name": "apple", "price": 200}'

curl -X POST \
     -H "Content-Type: application/json" \
     -d "$JSON_DATA" \
     "$API_ENDPOINT"
