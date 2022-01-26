#!/bin/sh
(cd backend && npx dotenv -e ../.env -e .env go run .)