FROM golang:bullseye

# install dependencies
RUN apt clean && apt-get update && apt-get -y install alien

WORKDIR /client

COPY client ./

RUN ls /client

WORKDIR /service

COPY service ./

RUN go mod tidy

RUN ls /service

RUN go build -o /service/main main.go

RUN mkdir -p data

# Add additional dependencies below ...

ENTRYPOINT [ "/service/main" ]