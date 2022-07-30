FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /docker-btc-service

EXPOSE 7777

CMD [ "/docker-btc-service" ]