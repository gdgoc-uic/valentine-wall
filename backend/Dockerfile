FROM golang:1.19-bullseye AS builder

WORKDIR /app

ENV CGOENABLED=1

COPY ./go.mod .
COPY ./go.sum .
COPY ./backend/*.go ./backend/
COPY ./backend/migrations ./backend/migrations
COPY ./backend/models ./backend/models
COPY ./backend/renderer_assets ./backend/renderer_assets
COPY ./backend/templates ./backend/templates
COPY ./backend/_data ./backend/_data

# Download deps and build
RUN go mod download
RUN cd backend && go build -o /app/backend-server

# FROM scratch
FROM debian:bullseye-slim

WORKDIR /app

RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates

RUN update-ca-certificates

COPY --from=builder /app/backend-server .

EXPOSE 8090
ENTRYPOINT ["./backend-server", "serve", "--http=0.0.0.0:8090", "--dir=/app/pb_data"]
