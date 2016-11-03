FROM tianon/true

ENV BIND 0.0.0.0
ENV PORT 8080

ADD cmd/exago/exago /
CMD ["/exago"]
