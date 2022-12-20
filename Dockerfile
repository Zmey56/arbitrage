FROM golang:latest

ENV GO111MODULE=on
ENV TELEGRAM_BOT_TOKEN="5763797414:AAHJ8exgiqxHuW44SyEr15fKsWKPixNofVg"

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build -o main .

CMD ["/app/main"]