FROM golang:1.16-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /app

COPY . .

RUN go build -o main .

EXPOSE ${APP_PORT}

CMD ["./main"]