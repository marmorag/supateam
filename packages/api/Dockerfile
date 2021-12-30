FROM golang:1.16-alpine as builder
# passed by buildkit
ARG TARGETOS
ARG TARGETARCH

RUN apk --update add ca-certificates make git

RUN mkdir -p /app
WORKDIR /app

COPY ../../. ./
RUN cd packages/api && \
    make setup && \
    TARGETOS=${TARGETOS} TARGETARCH=${TARGETARCH} CGO_ENABLED=0 make build

FROM scratch
LABEL org.opencontainers.image.source=https://github.com/marmorag/supateam

COPY --from=builder /app/packages/api/supateam /supateam
COPY --from=builder /app/packages/api/data.yaml.enc /data.yaml.enc

ENTRYPOINT ["/supateam"]
CMD ["serve"]