FROM golang:1.13 as build-env

WORKDIR /exago

COPY ../../go.mod .
COPY ../../go.sum .

FROM build-env
COPY ../.. .
WORKDIR /exago/src/consumer
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -a -tags netgo -ldflags '-w -extldflags "-static"' ./

ENTRYPOINT ["/go/bin/consumer"]