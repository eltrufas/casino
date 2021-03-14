FROM golang:alpine as builder
run apk add --update git

ENV CGO_ENABLED=0

WORKDIR /src

ADD go.mod go.sum ./
RUN go mod download

ADD . ./
RUN go build -o casino .


FROM scratch
COPY --from=builder /src/casino /
WORKDIR /
CMD ["./casino"]
