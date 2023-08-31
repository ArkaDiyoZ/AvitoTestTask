FROM golang:1.20

ENV GO111MODULE=on

WORKDIR /service

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /service/cmd

RUN go build -o main .

CMD ["./main"]