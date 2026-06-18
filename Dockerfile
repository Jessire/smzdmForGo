FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git
ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org,direct

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/smzdmPusher .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /opt/go
COPY --from=builder /out/smzdmPusher ./smzdmPusher
COPY config ./config
COPY template ./template
COPY data ./data

EXPOSE 9090
CMD ["./smzdmPusher"]
