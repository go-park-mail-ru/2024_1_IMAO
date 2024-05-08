FROM golang:1.22.0 as build

WORKDIR /var/backend

COPY cmd cmd
COPY internal internal
COPY uploads uploads

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod tidy \
    && go build -o main ./cmd/app/main.go

FROM alpine:edge as prod
RUN apk add bash

ENV DATABASE_URL=postgres://postgres:postgres@localhost:5432/IMAO_VOL4OK_2024

WORKDIR /root
COPY --from=build --chown=root:root /var/backend/main .
# RUN chmod 777 ./main

EXPOSE 8080
ENTRYPOINT ["${PWD}/main"]
