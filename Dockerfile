FROM golang:1.17-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -a -installsuffix cgo -o app .

FROM alpine:3.15

WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /build/app .

CMD ["./app"]
