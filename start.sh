#!/bin/bash

#Example of usage
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=discogs
export DB_SSL_MODE=disable
export DB_USERNAME=user
export DB_PASSWORD=password

go build
./discogs-parser -filename ./data_source/labels.xml -writer-type postgres

# Or for JSON output file saved into ./data_source/labels.json
./discogs-parser -filename ./data_source/labels.xml