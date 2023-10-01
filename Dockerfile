FROM golang:1.21.1-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/app

RUN CGO_ENABLED=0 GOOS=linux go build -o main

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/app/main /main
COPY --from=build-stage /app/app/.env /.env

ENTRYPOINT ["/main"]