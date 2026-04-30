FROM cgr.dev/chainguard/go:latest AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -tags "nosql gin" -ldflags="-s -w" -o /out/api ./cmd/api

FROM cgr.dev/chainguard/static:latest
COPY --from=builder /out/api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
