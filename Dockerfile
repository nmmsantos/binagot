FROM golang:1.17-bullseye AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN set -x && go mod download

COPY . .

RUN set -x && make

FROM gcr.io/distroless/base-debian11

COPY --from=build /app/bin/binagot /binagot

USER nonroot:nonroot

ENTRYPOINT ["/binagot"]
