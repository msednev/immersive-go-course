FROM golang:1.21 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY *.go ./
RUN GGO_ENABLED=0 GOOS=linux go build -o /docker-cloud

FROM gcr.io/distroless/base-debian12 AS release

WORKDIR /

COPY --from=build /docker-cloud /docker-cloud

EXPOSE 80

ENTRYPOINT [ "/docker-cloud" ]