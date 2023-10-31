FROM golang:1.21.1-alpine AS build-stage

WORKDIR /fill

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /fill/fill

RUN go build -o fill_db

FROM gcr.io/distroless/base-debian11 AS fill_db_stage

WORKDIR /

COPY --from=build-stage /fill/fill/fill_db /fill
COPY --from=build-stage /fill/fill/Netflix_Dataset.csv /Netflix_Dataset.csv

ENTRYPOINT ["/fill"]
