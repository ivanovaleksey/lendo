FROM golang:1.16

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY Makefile Makefile
COPY registry/ registry/
COPY pkg/ pkg/

ENV COMPONENT registry
RUN make build-app

CMD ["bin/registry"]
