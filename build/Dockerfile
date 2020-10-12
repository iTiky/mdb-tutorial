FROM golang:1.14-stretch

RUN apk update \
    && apk upgrade \
    && apk add --no-cache ca-certificates \
    && update-ca-certificates 2>/dev/null || true

COPY mdb-tutorial /
COPY resources/config.yml /etc/config.yml
COPY resources/csv /etc/csv

RUN chmod u+x /mdb-tutorial

EXPOSE 2412:2412

CMD ["/mdb-tutorial start --config /etc/config.yml"]
