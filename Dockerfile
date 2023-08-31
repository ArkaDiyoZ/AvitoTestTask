FROM golang:1.20 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main /app/cmd/main/main.go

FROM debian:buster-slim

COPY --from=build /app/main /app/main

CMD ["/app/main"]