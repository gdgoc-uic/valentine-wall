#!/bin/sh
if [ -x "$(command -v docker-compose)" ]; then
  docker-compose -f docker-compose.yml -f docker-compose.prod.yml up --no-deps --build -d $@
else
  docker compose -f docker-compose.yml -f docker-compose.prod.yml up --no-deps --build -d $@
fi