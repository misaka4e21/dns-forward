FROM golang:1.19-alpine AS build

WORKDIR /src

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /dns-forward

FROM alpine

COPY --from=build /dns-forward /dns-forward
ENTRYPOINT [ "/dns-forward" ]

ENV UPSTREAM_DNS_IP 127.0.0.11
ENV UPSTREAM_DNS_PORT 53
ENV LOCAL_DNS_PORT 53
EXPOSE 53/udp