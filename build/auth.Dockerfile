FROM golang:1.22.0 as build

WORKDIR /var/backend

COPY cmd cmd
COPY internal internal
COPY uploads uploads
COPY .env .env

COPY go.mod .
COPY go.sum .


RUN go mod tidy
RUN go build -o main ./cmd/auth_service/main.go
COPY --from=build /var/backend/.env ./.env

EXPOSE 8081
EXPOSE 7071

CMD ["./main"]
