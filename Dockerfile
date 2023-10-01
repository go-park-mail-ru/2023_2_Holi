FROM golang:1.21.1-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/app

#RUN go build -o main .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

CMD ["./main"]

