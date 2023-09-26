FROM golang:1.20 as build-env

RUN mkdir -p /opt/system-monitoring
WORKDIR /opt/system-monitoring

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build \
        -o /opt/system-monitoring/service \
        cmd/monitoring/main.go

FROM alpine:latest

COPY --from=build-env /opt/system-monitoring/service /opt/system-monitoring/service
COPY ./configs/config.toml /etc/system-monitoring/config.toml

ENTRYPOINT ["/opt/system-monitoring/service", "-config", "/etc/system-monitoring/config.toml"]