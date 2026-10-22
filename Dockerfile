FROM golang:1.26.10-alpine3.23 AS builder

WORKDIR /src

COPY . .
RUN CGO_ENABLED=0 go build -o /domain-exporter ./cmd/domain-exporter

FROM alpine:3.24

RUN apk add --no-cache ca-certificates-bundle=20261020-r0

COPY --from=builder /domain-exporter /bin/domain-exporter

EXPOSE      9553
USER        nobody
ENTRYPOINT  [ "/bin/domain-exporter" ]
CMD         [ "--config_path=/etc/domain-exporter/config.yaml" ]
