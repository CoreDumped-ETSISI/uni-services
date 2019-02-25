FROM golang:alpine AS build-env

# Install certs
RUN apk add --no-cache ca-certificates

ADD . /go/src/github.com/guad/status

RUN cd /go/src/github.com/guad/status && go build -o app


FROM alpine

# Import certificates
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs

WORKDIR /app
COPY --from=build-env /go/src/github.com/guad/status/app /app/
ENTRYPOINT ./app