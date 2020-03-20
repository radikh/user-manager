#!/bin/bash

curl --request PUT --data 5432 localhost:8500/v1/kv/db-postgres
curl --request PUT --data 9000 localhost:8500/v1/kv/logger-graylog
