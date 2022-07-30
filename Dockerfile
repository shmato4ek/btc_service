FROM node:12-alpine
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN yarn install --production
CMD ["go", "main.go"]
EXPOSE 8080