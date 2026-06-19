FROM golang:1.26.4-alpine AS builder

WORKDIR /src

COPY . .
RUN CGO_ENABLED=0 go build -o /domain-exporter ./cmd/domain-exporter

FROM alpine:3.24

COPY --from=builder /domain-exporter /bin/domain-exporter

EXPOSE      9553
USER        nobody
ENTRYPOINT  [ "/bin/domain-exporter" ]
CMD         [ "--config_path=/etc/domain-exporter/config.yaml" ]
