FROM golang:1.17.6 AS builder

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .
COPY ./backend .
COPY ./postal_office/types ./postal_office/types

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

EXPOSE 4000
ENTRYPOINT ["./backend-server"]