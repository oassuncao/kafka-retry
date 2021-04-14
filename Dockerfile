FROM golang:1.14 AS builder
WORKDIR /build

# Resolve dependencies first
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -o kafka-retry cmd/retry/retry.go

FROM ubuntu:bionic AS final
RUN apt-get update -y
RUN apt-get install wget gnupg -y

RUN wget -qO - https://packages.confluent.io/deb/5.5/archive.key | apt-key add -
RUN echo "deb [arch=amd64] https://packages.confluent.io/deb/5.5 stable main" > /etc/apt/sources.list.d/confluent.list
RUN apt-get update -y
RUN apt-get install librdkafka-dev -y

WORKDIR /app
COPY --from=builder /build/kafka-retry .

CMD ["./kafka-retry"]
