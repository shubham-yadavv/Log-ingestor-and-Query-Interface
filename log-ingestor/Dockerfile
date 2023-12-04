FROM golang:alpine AS builder

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /go/bin/app ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /go/bin/app .

EXPOSE 3000

CMD ["./app"]