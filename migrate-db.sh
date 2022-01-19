#!/bin/sh

# do not execute this - just for sample

npx dotenv -e .env -- sql-migrate up -config ./backend/dbconfig-nondocker.yml