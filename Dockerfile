FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/github.com/shauncampbell/namer/
COPY . .

RUN go build -o /go/bin/dapper github.com/shauncampbell/namer/cmd/namer

FROM alpine:3.12

COPY --from=builder /go/bin/namer /go/bin/namer
LABEL maintainer="Shaun Campbell <docker@shaun.scot>"
COPY docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh

VOLUME /config.yaml
ENV LDAP_BASE "dc=home,dc=lab"
EXPOSE 53

ENTRYPOINT ["./docker-entrypoint.sh"]
