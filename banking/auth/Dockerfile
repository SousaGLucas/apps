FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY app ./app
COPY domain ./domain
COPY gateways ./gateways
COPY telemetry ./telemetry

RUN CGO_ENABLED=0 GOOS=linux go build -v -o build/bin ./app/service

FROM ubuntu:20.04 as runner

RUN apt-get update -y && apt-get install -y curl

RUN useradd --create-home --uid 1000 app
WORKDIR /home/app
USER 1000

COPY --from=builder /app/build/bin /build/
ENTRYPOINT ["/build/bin"]
