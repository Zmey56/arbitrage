FROM multiarch/qemu-user-static:x86_64-xtensaeb-7.2.0-1 as builder

ENV GO_VERSION=1.15
RUN apk add --no-cache curl \
    && curl -fsSL "https://golang.org/dl/go$GO_VERSION.src.tar.gz" -o go.tar.gz \
    && tar -C /usr/local -xzf go.tar.gz \
    && rm go.tar.gz \
    && cd /usr/local/go/src \
    && ./make.bash

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

COPY . $GOPATH/src/myapp/

RUN cd $GOPATH/src/myapp/cmd/telegramfiats \
    && CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o main .

FROM alpine:latest
COPY --from=builder $GOPATH/src/myapp/cmd/telegramfiats/main .
CMD ["./main"]

