FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc git musl-dev sqlite-dev
ENV CGO_ENABLED=1
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/smzdmPusher .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates sqlite-libs tzdata

WORKDIR /opt/go
COPY --from=builder /out/smzdmPusher ./smzdmPusher
COPY config ./config
COPY template ./template
COPY data ./data

EXPOSE 9090
CMD ["./smzdmPusher"]
