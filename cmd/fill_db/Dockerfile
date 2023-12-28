FROM golang:1.21.1-alpine AS build-stage

WORKDIR /fill

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /fill/cmd/fill_db

RUN go build -o fill_db

FROM gcr.io/distroless/base-debian11 AS fill_db_stage

WORKDIR /

COPY --from=build-stage /fill/cmd/fill_db/fill_db /fill_db
COPY --from=build-stage /fill/cmd/fill_db/Netflix_Dataset.csv /Netflix_Dataset.csv
COPY --from=build-stage /fill/cmd/fill_db/Actors.csv /Actors.csv

ENTRYPOINT ["/fill_db"]
