FROM golang:1.21-alpine3.18 AS builder

WORKDIR /root/sti
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -trimpath -o out/sti ./cmd/sti/*.go


FROM alpine:3.18

RUN apk add \
  ca-certificates \
  iperf3
COPY --from=builder /root/sti/out/sti /usr/local/bin/sti

ENTRYPOINT ["sti"]
