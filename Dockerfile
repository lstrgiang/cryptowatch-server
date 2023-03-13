FROM golang:1.18-buster AS build

ARG GITHUB_TOKEN

WORKDIR /go/src/crypto-watch/server

RUN git config --global url.https://$GITHUB_TOKEN@github.com/.insteadOf https://github.com/

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/app .

FROM alpine:3.11

ENV GIN_MODE release

COPY --from=0 /usr/local/bin/app /usr/local/bin/app
COPY migrations migrations

RUN apk add --no-cache ca-certificates

ENTRYPOINT ["app", "server"]
