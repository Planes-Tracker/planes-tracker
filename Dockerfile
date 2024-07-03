ARG GOLANG_VERSION=1.20

FROM golang:${GOLANG_VERSION}-alpine as build

WORKDIR /src

COPY go.mod .
COPY go.sum .
COPY main.go .
COPY ./app app/
COPY ./internal internal/

RUN set -eux; \
    apk add musl-dev gcc

RUN set -eux; \
    go mod download

RUN set -eux; \
    env CGO_ENABLED=1 go build -ldflags "-s -w" -o /bin/planes-tracker



FROM alpine

WORKDIR /app

COPY config.json .
COPY --from=build /bin/planes-tracker /bin/planes-tracker

RUN set -eux; \
    apk add libc6-compat

CMD [ "/bin/planes-tracker" ]
