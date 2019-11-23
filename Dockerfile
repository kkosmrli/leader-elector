# Build stage
FROM golang:alpine3.10 AS builder
WORKDIR /src
ADD . /src
RUN cd /src && go build -o elector cmd/leader-elector/main.go

FROM alpine:3.10.3
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /src/elector /app/
ENTRYPOINT ["./elector"]