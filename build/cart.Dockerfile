FROM golang:1.22.0 as build

WORKDIR /var/backend

COPY cmd cmd
COPY internal internal
COPY uploads uploads

COPY go.mod .
COPY go.sum .


RUN go mod tidy
RUN go build -o main ./cmd/cart_service/main.go

EXPOSE 8083

CMD ["./main"]
