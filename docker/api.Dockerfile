FROM golang:1.16-alpine

RUN apk add --update make

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY Makefile Makefile
COPY api/ api/
COPY pkg/ pkg/

ENV COMPONENT api
RUN make build-app

EXPOSE 8000
CMD ["bin/api"]
