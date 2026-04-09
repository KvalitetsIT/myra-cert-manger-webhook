FROM --platform=$BUILDPLATFORM golang:1.25.9-alpine AS builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
ARG TARGETOS TARGETARCH
RUN --mount=type=cache,id=gomod-${TARGETOS}-${TARGETARCH},target=/go/pkg/mod \
    go mod download
COPY . .
RUN --mount=type=cache,id=gomod-${TARGETOS}-${TARGETARCH},target=/go/pkg/mod \
    --mount=type=cache,id=gobuild-${TARGETOS}-${TARGETARCH},target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags="-s -w" -o /app ./cmd/app/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app/app
USER 65532:65532
ENTRYPOINT ["/app/app"]
