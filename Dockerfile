FROM alpine:3.12.1 AS prerequisite

###

FROM golang:1.16.3 AS base

WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download

###

FROM base AS build

COPY . .
RUN make build

###

FROM prerequisite

COPY --from=build /src/grafana2vonage /

ENTRYPOINT ["/grafana2vonage"]
