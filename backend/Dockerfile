FROM golang:1.17.6-alpine

ARG env
ARG database_prefix
ARG firebase_config

ENV ENV=${env}
ENV DATABASE_PREFIX=${database_prefix}
ENV FIREBASE_CONFIG=${firebase_config}

WORKDIR /app

RUN apk add build-base sqlite

COPY ./go.mod ./
COPY ./go.sum ./
COPY ./backend .

# Download deps
RUN go mod download
RUN go install github.com/rubenv/sql-migrate/sql-migrate@latest

# Run migrations
RUN touch ${DATABASE_PREFIX}_${ENV}.db
RUN sql-migrate up -env ${ENV}

# Build
COPY ${firebase_config} .
RUN ls
RUN go build -o /backend-server

EXPOSE 4000

ENTRYPOINT ["/backend-server"]