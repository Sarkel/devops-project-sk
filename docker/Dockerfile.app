FROM golang:1.25.5-alpine3.23 AS base

WORKDIR /app

RUN apk add --no-cache make

COPY Makefile .
COPY app ./app
COPY common ./common
COPY seeder ./seeder
COPY go.work .

FROM base AS build
RUN make build-app

FROM base AS test

RUN make test-app

FROM scratch AS artifact-export

COPY --from=test /app/app/coverage.out /coverage.out

FROM alpine:3.23 AS runtime

WORKDIR /app

# install common libraries/os packages

FROM runtime AS api

COPY --from=build /app/app/bin/api /app/api

ENTRYPOINT ["/app/api"]

FROM runtime AS crawler

COPY --from=build /app/app/bin/crawler /app/crawler

ENTRYPOINT ["/app/crawler"]

FROM runtime AS reader

COPY --from=build /app/app/bin/reader /app/reader

ENTRYPOINT ["/app/reader"]
