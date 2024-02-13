#!/bin/sh
cd backend && (dotenv -e ../.env -e .env -- go run . $@)
