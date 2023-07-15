FROM golang:1.20-alpine

ARG UID=1000
ARG GID=1000

RUN set -ex \
  && addgroup -g ${GID} pubgo \
  && adduser -u ${UID} -h /opt/pubgo -s /bin/sh -G pubgo -D pubgo

ADD . /srv/
WORKDIR /srv/

ENV GO111MODULE="on"

RUN go mod tidy && go build -o main *.go

ENTRYPOINT su -c "/srv/main -config /opt/pubgo/config.yaml -content_dir /opt/pubgo" pubgo

EXPOSE 8080
