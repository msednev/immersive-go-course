FROM golang:1.21 as build

WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY consumer ./consumer
COPY producer ./producer
COPY utils ./utils

RUN mkdir -p /out
RUN CGO_ENABLED=1 go build -o /out ./...

FROM gcr.io/distroless/cc-debian12 as producer

WORKDIR /

COPY --from=build /out/producer /
COPY cron.txt /

ENTRYPOINT ["/producer", "/cron.txt"]

FROM debian:bookworm-slim as consumer

WORKDIR /

COPY --from=build /out/consumer /

ENTRYPOINT ["/consumer"]