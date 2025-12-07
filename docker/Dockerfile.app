FROM golang:1.25.5-alpine3.23 as base

WORKDIR /app

RUN apk add --no-cache make

COPY Makefile .
COPY app ./app

FROM base as build
RUN make build-app

FROM base as test

FROM alpine:3.23 as runtime

WORKDIR /app

# install common libraries/os packages

FROM runtime as api

COPY --from=build /app/app/bin/api /app/api

CMD ["/app/api"]

FROM runtime as carler

COPY --from=build /app/app/bin/crawler /app/crawler

CMD ["/app/crawler"]

FROM runtime as reader

COPY --from=build /app/app/bin/reader /app/reader

CMD ["/app/reader"]
