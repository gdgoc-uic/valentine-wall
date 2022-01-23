#!/bin/sh

docker-compose build $@
docker-compose up --no-deps $@