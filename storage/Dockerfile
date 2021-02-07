FROM golang:1.15-alpine3.12

RUN mkdir /app
ADD . /app

WORKDIR /app

RUN go build -o storage

CMD [ "/app/storage" ]