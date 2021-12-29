FROM golang:1.16-alpine as builder
# passed by buildkit
ARG TARGETOS
ARG TARGETARCH

RUN apk --update add ca-certificates make git

RUN mkdir -p /app
WORKDIR /app

COPY . .
RUN TARGETOS=${TARGETOS} TARGETARCH=${TARGETARCH} CGO_ENABLED=0 make build

FROM scratch

COPY --from=builder /app/supateam /supateam

ENTRYPOINT ["/supateam"]
CMD ["serve"]