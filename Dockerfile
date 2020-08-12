FROM golang:1.14.6 AS builder

ENV GO111MODULE=on
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
EXPOSE 8080
RUN GOOS=linux go build -o workout github.com/michael-diggin/workout/server

ENTRYPOINT ["./workout"]