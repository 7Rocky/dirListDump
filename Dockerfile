FROM golang:1.16-alpine AS build

WORKDIR /build
RUN apk update && apk add upx

COPY go.* .
COPY *.go .

RUN go mod download && go mod tidy
RUN go build --ldflags='-s -w' -o dldump *.go
RUN upx --ultra-brute dldump


FROM alpine

RUN addgroup -S gopher
RUN adduser --disabled-password --gecos '' --home /home/gopher --ingroup gopher --shell /bin/sh gopher

USER gopher
WORKDIR /home/gopher

COPY --from=build --chown=gopher /build/dldump /usr/local/bin

CMD [ "dldump" ]
