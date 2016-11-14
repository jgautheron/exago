FROM alpine:3.4

RUN apk add --no-cache ca-certificates curl

ENV BIND 0.0.0.0
ENV PORT 8080

ADD cmd/exago/exago /
ENTRYPOINT ["/exago"]
