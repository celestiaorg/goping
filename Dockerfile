# FROM golang:alpine3.15 AS development
FROM docker.io/golang:1.21-alpine3.18 AS development
ARG arch=x86_64

# ENV CGO_ENABLED=0
WORKDIR /go/src/app/
COPY . /go/src/app/

RUN apk update && apk add --no-cache \
    make && \
    mkdir -p /build/ && \
    make build && \
    cp ./bin/* /build/

#----------------------------#

FROM docker.io/alpine:3.18.4 AS production

WORKDIR /app/
COPY --from=development /build .

ENV PATH=$PATH:/app/
ENV SERVE_ADDR="0.0.0.0:8080"
CMD ["sh", "-c", "goping serve --serve-addr ${SERVE_ADDR}"]

