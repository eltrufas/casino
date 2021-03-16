FROM golang:alpine as builder
RUN apk add git build-base tzdata ca-certificates

WORKDIR /src

ADD go.mod go.sum ./
RUN go mod download

ADD . ./
RUN go build -ldflags="-extldflags=-static" -o casino .

FROM scratch
COPY --from=builder /etc/ssl/ /etc/ssl/certs
COPY --from=builder /src/casino /
WORKDIR /
ENTRYPOINT ["./casino"]
