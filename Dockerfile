FROM golang:alpine

RUN mkdir /app

WORKDIR /app

ADD go.mod .
ADD go.sum .

RUN go mod download
ADD . .

EXPOSE 8080
CMD go run main.go && go run receive.go
