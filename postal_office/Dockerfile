FROM golang:1.17.6 AS builder
WORKDIR /app
COPY ./go.mod .
COPY ./go.sum .
COPY ./postal_office ./postal_office
RUN go mod download
RUN go build -o ./postal-office ./postal_office

FROM debian:buster-slim
WORKDIR /app
RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates \
 && update-ca-certificates
COPY --from=builder /app/postal-office .
EXPOSE 3350
ENTRYPOINT ["./postal-office"]