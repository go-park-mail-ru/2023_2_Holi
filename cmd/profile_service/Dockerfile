FROM golang:1.21.1-alpine AS build-stage

WORKDIR /app

COPY ../../go.mod ../../go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/profile_service

RUN go build -o main

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/cmd/profile_service/main /main

ENTRYPOINT ["/main"]
