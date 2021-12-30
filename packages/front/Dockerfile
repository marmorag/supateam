FROM node:14 as app-builder
WORKDIR /app

COPY . .
RUN cd packages/front && \
    yarn && \
    yarn build

FROM alpine:3.13.2 AS thttpd-builder
ARG THTTPD_VERSION=2.29

RUN apk add gcc musl-dev make
RUN wget https://www.acme.com/software/thttpd/thttpd-${THTTPD_VERSION}.tar.gz \
  && tar xzf thttpd-${THTTPD_VERSION}.tar.gz \
  && mv /thttpd-${THTTPD_VERSION} /thttpd
RUN cd /thttpd \
  && ./configure \
  && make CCOPT='-O2 -s -static' thttpd

RUN adduser -D static-user

FROM scratch as application
LABEL org.opencontainers.image.source=https://github.com/marmorag/supateam

COPY --from=thttpd-builder /etc/passwd /etc/passwd
COPY --from=thttpd-builder /thttpd/thttpd /bin/thttpd

USER static-user
COPY --from=app-builder /app/packages/front/dist /static

EXPOSE 8000
ENTRYPOINT ["/bin/thttpd"]
CMD ["-D", "-h", "0.0.0.0", "-p", "8000", "-d", "/static", "-u", "static-user", "-l", "-", "-M", "60"]