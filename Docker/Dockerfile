FROM golang:latest as build
WORKDIR /build
COPY . /build
RUN go build .

FROM grafana/promtail:2.1.0

# Copy the mock-Loki
COPY --from=build /build/promtail-debug /usr/bin/promtail-debug

# Copy the base config
COPY Docker/config.yaml /config.yaml

# Expose all the ports that are useful
EXPOSE 9007
ENTRYPOINT ["/usr/bin/promtail-debug"]