FROM golang:1.19 AS builder

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .
COPY ./backend .
COPY ./backend/migrations ./backend/migrations
COPY ./backend/models ./backend/models

# Download deps and build
RUN go mod download
RUN go build -o ./backend-server

# FROM scratch
FROM debian:buster-slim

WORKDIR /app

RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates

RUN update-ca-certificates

COPY --from=builder /app/backend-server .

EXPOSE 8090
ENTRYPOINT ["./backend-server", "serve", "--http=0.0.0.0:8090", "--dir=/app/pb_data"]