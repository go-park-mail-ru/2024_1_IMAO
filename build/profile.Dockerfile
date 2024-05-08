FROM golang:1.22.0 as build

WORKDIR /var/backend

COPY cmd cmd
COPY internal internal
COPY uploads uploads

COPY go.mod .
COPY go.sum .


RUN go mod tidy
RUN go build -o main ./cmd/profile_service/main.go

EXPOSE 8082
EXPOSE 7072

CMD ["./main"]
