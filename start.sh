#!/bin/bash

#Example of usage
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=discogs
export DB_SSL_MODE=disable
export DB_USERNAME=user
export DB_PASSWORD=password

go build
./discogs -filename ./source/labels.xml -block-size 1000 -block-skip 0 -block-limit 2147483647 -writer-type postgres