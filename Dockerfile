FROM alpine:3.4

RUN apk add --no-cache ca-certificates curl

ENV BIND 0.0.0.0
ENV PORT 8080

HEALTHCHECK --interval=1m --timeout=2s \
  CMD curl -f http://localhost:8080/_health || exit 1

ADD cmd/exago/exago /

EXPOSE 8080
ENTRYPOINT ["/exago"]
