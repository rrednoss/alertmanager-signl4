FROM golang:1.17-bullseye as builder
WORKDIR /app
COPY . /app/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -o main

FROM scratch
LABEL org.opencontainers.image.authors="Rene Rednoß"
LABEL org.opencontainers.image.description="Transforms the Alertmanager payload to a more SIGNL4 friendly one."
LABEL org.opencontainers.image.licenses="MIT License"
LABEL org.opencontainers.image.url="https://github.com/rrednoss/alertmanager-signl4"
LABEL org.opencontainers.image.version="v0.1.0"
COPY --from=builder /app/main .
EXPOSE 9094/tcp
ENTRYPOINT [ "/main" ]