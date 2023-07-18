FROM golang:1.19-alpine AS builder

WORKDIR /build

RUN apk update && apk upgrade
RUN apk add --no-cache postgresql-dev build-base postgresql-client

COPY . .

RUN go build -o app github.com/profectus200/contact-book-bot/cmd/bot

EXPOSE 8080

CMD ["./app"]
