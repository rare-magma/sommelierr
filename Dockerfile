FROM docker.io/library/golang:alpine AS builder
WORKDIR /app
ENV CGO_ENABLED=0
COPY go.mod go.mod
COPY internal internal
COPY cmd cmd
RUN go build -ldflags "-s -w" -trimpath -o app ./cmd/server

FROM cgr.dev/chainguard/static:latest
COPY --from=builder /app/app /usr/bin/app
ENTRYPOINT ["/usr/bin/app"]