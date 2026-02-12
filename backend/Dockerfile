FROM golang:1.19-bullseye AS builder

WORKDIR /app

ENV CGOENABLED=1

COPY ./go.mod .
COPY ./go.sum .
COPY ./backend/*.go ./
COPY ./backend/migrations ./migrations
COPY ./backend/models ./models
COPY ./backend/renderer_assets ./renderer_assets
COPY ./backend/templates ./templates
COPY ./backend/_data ./_data

# Download deps and build
RUN go mod download
RUN go build -o ./backend-server

# FROM scratch
FROM debian:bullseye-slim

WORKDIR /app

RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates

RUN update-ca-certificates

COPY --from=builder /app/backend-server .

EXPOSE 8090
ENTRYPOINT ["./backend-server", "serve", "--http=0.0.0.0:8090", "--dir=/app/pb_data"]
