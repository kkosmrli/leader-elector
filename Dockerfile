ARG ARCH=amd64

# Build stage
FROM golang:1.16-alpine3.13 AS builder
WORKDIR /src
ADD . /src
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} go get github.com/alexflint/go-arg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} go build -ldflags='-s -w -extldflags "-static"' -o elector cmd/leader-elector/main.go

FROM gcr.io/distroless/static:nonroot-${ARCH}
USER nonroot:nonroot
WORKDIR /app
COPY --from=builder --chown=nonroot:nonroot /src/elector /app/
ENTRYPOINT ["./elector"]
